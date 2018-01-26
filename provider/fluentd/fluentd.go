package fluentd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/rancher/logging-controller/provider"
	"github.com/rancher/types/apis/logging.cattle.io/v3"
	"github.com/rancher/types/config"
)

const (
	fluentdLogFile      = "fluentd.log"
	fluentdPidFile      = "fluentd.pid"
	fluentdConfigFile   = "fluentd.conf"
	clusterConfigFile   = "cluster.conf"
	projectConfigFile   = "project.conf"
	serviceConfigFile   = "service.conf"
	clusterTemplateFile = "cluster_template.conf"
	projectTemplateFile = "project_template.conf"
	serviceTemplateFile = "service_template.conf"
)

var (
	fluentdProcess *exec.Cmd
	fluentdTimeout = 1 * time.Minute
)

type Provider struct {
	configDir     string
	stopCh        chan struct{}
	dryRun        bool
	startCmd      string
	pidPath       string
	logPath       string
	cfgPath       string
	ClusterConfig ClusterConfig
	ProjectConfig ProjectConfig
}

type ClusterConfig struct {
	logLister v3.LoggingLister
	configDir string
}

type ProjectConfig struct {
	pjlogLister v3.ProjectLoggingLister
	configDir   string
}

func (logp *Provider) Init(c *cli.Context, lg *config.LoggingContext) {
	logp.configDir = c.String("fluentd-config-dir")
	logp.dryRun = c.Bool("fluentd-dry-run")
	logp.logPath = path.Join(logp.configDir, fluentdLogFile)
	logp.cfgPath = path.Join(logp.configDir, fluentdConfigFile)
	logp.pidPath = path.Join(logp.configDir, fluentdPidFile)
	logp.startCmd = "fluentd " + "-c " + logp.cfgPath + " -d " + logp.pidPath + " --log " + logp.logPath
	logp.ClusterConfig = ClusterConfig{
		logLister: lg.Logging.Loggings("").Controller().Lister(),
	}
	logp.ProjectConfig = ProjectConfig{
		pjlogLister: lg.Logging.ProjectLoggings("").Controller().Lister(),
	}
}

func (cfg *ClusterConfig) Update() error {
	conf := make(map[string]interface{})
	loggings, err := cfg.logLister.List("", labels.NewSelector())
	if err != nil {
		return err
	}
	if len(loggings) == 0 {
		return fmt.Errorf("no resource logging exist")
	}
	conf["clusterTarget"] = loggings[0].Spec
	return update("cluster", conf)
}

func (cfg *ProjectConfig) Update() error {
	conf := make(map[string]interface{})
	pjloggings, err := cfg.pjlogLister.List("", labels.NewSelector())
	if err != nil {
		return err
	}
	if len(pjloggings) == 0 {
		return fmt.Errorf("no resource projectlogging exist")
	}

	var pjlogsTgs []v3.ProjectLoggingSpec
	for _, v := range pjloggings {
		pjlogsTgs = append(pjlogsTgs, v.Spec)
	}
	conf["projectTargets"] = pjlogsTgs
	return update("project", cfg.configDir, conf)
}

func update(logType, configDir string, conf map[string]interface{}) error {
	var w io.Writer
	tmpPath := path.Join(configDir, logType, ".tmp")
	w, err := os.Create(tmpPath)
	if err != nil {
		return errors.Wrap(err, "fail create create fluentd tmp config")
	}

	if _, err := os.Stat(tmpPath); err != nil {
		return errors.Wrap(err, "fail get created fluentd tmp config file")
	}
	templatePath := path.Join(configDir, logType+"_template.conf")
	var t *template.Template
	t, err = template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	// conf = getDate()
	err = t.Execute(w, conf)
	if err != nil {
		return err
	}

	// only change fluentd config when real change happen
	cfgPath := path.Join(configDir, logType+".conf")
	cfgEqual, err := isConfigEqual(cfgPath, tmpPath)
	if err != nil {
		return err
	}
	if cfgEqual {
		logrus.Info("config file not change, no need to reload")
		return nil
	}
	logrus.Info("config file changed, reloading")
	err = os.Rename(cfgPath, cfgPath+".bak")
	if err != nil {
		return errors.Wrap(err, "fail to rename config config file")
	}
	from, err := os.Open(tmpPath)
	if err != nil {
		return errors.Wrap(err, "fail to open tmp config file")
	}
	defer from.Close()

	to, err := os.OpenFile(cfgPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrap(err, "fail to open current config file")
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return errors.Wrap(err, "fail to copy config file")
	}
	if err = to.Sync(); err != nil {
		return errors.Wrap(err, "fail to sync config file")
	}
	return nil
}

func init() {
	logp := Provider{
		stopCh: make(chan struct{}),
	}
	provider.RegisterProvider(logp.GetName(), logp)
}

func (logp *Provider) GetName() string {
	return "fluentd"
}

// func (logp *Provider) Run() {
// 	if logp.dryRun {
// 		return
// 	}
// 	cfg, err := infraconfig.GetLoggingConfig(api.NamespaceAll, loggingv1.LoggingName)
// 	if err != nil {
// 		logrus.Errorf("fail get logging config, details: %s", err.Error())
// 		<-logp.stopCh
// 		return
// 	}
// 	if err = logp.cfg.write(cfg); err != nil {
// 		logrus.Errorf("fail write fluentd config, details: %s", err.Error())
// 		<-logp.stopCh
// 		return
// 	}

// 	if err := logp.StartFluentd(); err != nil {
// 		logrus.Errorf("fail start fluentd, details: %s", err.Error())
// 		<-logp.stopCh
// 		return
// 	}
// 	<-logp.stopCh
// }

func (logp *Provider) Stop() error {
	logrus.Warnf("shutting down provider %s", logp.GetName())
	close(logp.stopCh)
	return nil
}

func (logp *Provider) GetConfig(name string) provider.LogProviderConfig {
	if name == "cluster" {
		return logp.ClusterConfig
	} else if name == "project" {
		return logp.ProjectConfig
	}
	return nil
}

func (logp *Provider) Start() error {
	if logp.dryRun {
		return nil
	}
	cmd := exec.Command("sh", "-c", logp.startCmd)
	logrus.Infof("fluentd start command: %s", logp.startCmd)
	var buf bytes.Buffer
	cmd.Stdout = &buf

	cmd.Start()

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	timeout := time.After(fluentdTimeout)

	select {
	case <-timeout:
		cmd.Process.Kill()
		return errors.New("Fluentd command timed out")
	case err := <-done:
		logrus.Error("Fluentd Output:", buf.String())
		if err != nil {
			logrus.Error("Fluentd return a Non-zero exit code:", err)
			return err
		}
	}
	return nil
}

func (logp *Provider) Reload() error {
	if logp.dryRun {
		return nil
	}
	pidFile, err := ioutil.ReadFile(logp.pidPath)
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(pidFile)))
	if err != nil {
		return fmt.Errorf("fail parsing pid from %s: %s", pidFile, err)
	}

	if pid <= 0 {
		logrus.Warning("Fluentd not start yet, could not reload")
		return nil
	}
	if _, err := os.FindProcess(pid); err != nil {
		return fmt.Errorf("fail find process pid: %d, details: %v", pid, err)
	}

	if err = syscall.Kill(pid, syscall.SIGHUP); err != nil {
		return fmt.Errorf("fail reloading, details: %v", err)
	}
	return nil
}

func isConfigEqual(beforePath, afterPath string) (bool, error) {
	f1, err := ioutil.ReadFile(afterPath)

	if err != nil {
		return false, errors.Wrapf(err, "fail read file %s", afterPath)
	}

	f2, err := ioutil.ReadFile(beforePath)

	if err != nil {
		return false, errors.Wrapf(err, "fail read file %s", beforePath)
	}
	return bytes.Equal(beforePath, afterPath), nil
}
