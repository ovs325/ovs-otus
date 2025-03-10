package manager

import (
	cf "bruteforce/config"
	mk "bruteforce/internal/manager/mocks"
	tp "bruteforce/internal/types"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	once          sync.Once // Контролирует все ниже расположенные переменные
	Cfg           cf.Config
	Ctx           context.Context
	configContent string

	Manager     *GroupManager
	groupsMap   map[string]*mk.GroupItem
	specListMap map[string]*mk.SpecListsManager
)

func initGlobalVars(t *testing.T) (err error) {
	initV := func() {
		Ctx = context.Background()

		configContent = `
params:
  login:
    lenQueue: 10000
    capacity: 10
    leakageRate: 0.16666666666666666
  pass:
    lenQueue: 10000
    capacity: 10
    leakageRate: 1.6666666666666667
  ip:
    lenQueue: 10000
    capacity: 10
    leakageRate: 16.666666666666668
default:
  lenQueue: 10000
  capacity: 10
  leakageRate: 0.5
logger:
  level: info
http_server:
  scheme: http
  host: 127.0.0.1
  port: 3009
is_test: true
`

		tempDir := t.TempDir()

		configPathTempDir := filepath.Join(tempDir, "config.yaml")
		assert.Nil(t, os.WriteFile(configPathTempDir, []byte(configContent), 0644), "failed to write config file")

		// Загружаем конфигурацию
		Cfg, err = cf.LoadConfig(tempDir)
		assert.NoError(t, err)

		Manager = &GroupManager{
			cfg:         &Cfg,
			groupsMap:   groupMap{},
			specListMap: map[string]SpecListsManager{},
		}

		groupsMap = map[string]*mk.GroupItem{}
		specListMap = map[string]*mk.SpecListsManager{}

		for name := range Manager.cfg.Params {
			m := mk.NewGroupItem(t)
			groupsMap[name] = m
			Manager.groupsMap[name] = m
			spl := mk.NewSpecListsManager(t)
			specListMap[name] = spl
			Manager.specListMap[name] = spl
		}
	}
	once.Do(initV)
	return
}

func TestInit(t *testing.T) {
	assert.NoError(t, initGlobalVars(t))

	assert.NotNil(t, Ctx)

	assert.NotNil(t, configContent)
	assert.NotNil(t, Cfg)
	require.NotNil(t, Cfg.Params)

	require.NotNil(t, Manager)
	assert.Equal(t, len(Cfg.Params), len(Manager.groupsMap))
	assert.Equal(t, len(Cfg.Params), len(Manager.specListMap))

	assert.NotNil(t, groupsMap)
	assert.NotNil(t, specListMap)
	assert.Equal(t, len(Cfg.Params), len(groupsMap))
	assert.Equal(t, len(Cfg.Params), len(specListMap))
}

func TestNewManager(t *testing.T) {
	manager, err := NewManager("../test/")
	testErr := fmt.Errorf("Ошибка при загрузке конфига: path = %s, err - %w", "../test/", err)
	assert.ErrorIs(t, testErr, err)
	require.Nil(t, manager)

	tempDir := t.TempDir()
	configPathTempDir := filepath.Join(tempDir, "config.yaml")
	assert.Nil(t, os.WriteFile(configPathTempDir, []byte(configContent), 0644), "failed to write config file")
	manager, err = NewManager(tempDir)
	assert.NoError(t, err)
	require.NotNil(t, manager)

	for name, params := range manager.cfg.Params {
		groupParams := manager.GetParams(name)
		assert.Equal(t, params.Capacity, groupParams.Capacity)
		assert.Equal(t, params.LeakageRate, groupParams.LeakageRate)
		assert.Equal(t, params.LenQueue, groupParams.LenQueue)
	}

}

func TestNewGroup(t *testing.T) {
	assert.NoError(t, initGlobalVars(t))
	tempDir := t.TempDir()
	configPathTempDir := filepath.Join(tempDir, "config.yaml")
	assert.Nil(t, os.WriteFile(configPathTempDir, []byte(configContent), 0644), "failed to write config file")
	manager, err := NewManager(tempDir)
	assert.NoError(t, err)
	require.NotNil(t, manager)

	params := cf.GroupParams{LenQueue: 100, Capacity: 10, LeakageRate: 1}
	assert.NoError(t, manager.NewGroup("test", params))

	groupParams := manager.GetParams("test")
	assert.Equal(t, 4, len(manager.groupsMap))
	assert.Equal(t, int64(10), groupParams.Capacity)
	assert.Equal(t, 1.0, groupParams.LeakageRate)
	assert.Equal(t, int64(100), groupParams.LenQueue)
	assert.Equal(t, 4, len(manager.cfg.Params))
	assert.Equal(t, int64(10), manager.cfg.Params["test"].Capacity)
	assert.Equal(t, 1.0, manager.cfg.Params["test"].LeakageRate)
	assert.Equal(t, int64(100), manager.cfg.Params["test"].LenQueue)

	err = manager.NewGroup("ip", params)
	assert.Equal(t, err, fmt.Errorf("такая группа уже существует: name = ip"))

	manager.cfg = nil
	assert.NoError(t, manager.NewGroup("cfgNil", params))

	groupParams = manager.GetParams("cfgNil")
	assert.Equal(t, 5, len(manager.groupsMap))
	assert.Equal(t, int64(10), groupParams.Capacity)
	assert.Equal(t, 1.0, groupParams.LeakageRate)
	assert.Equal(t, int64(100), groupParams.LenQueue)
	assert.Equal(t, 1, len(manager.cfg.Params))
	assert.Equal(t, int64(10), manager.cfg.Params["cfgNil"].Capacity)
	assert.Equal(t, 1.0, manager.cfg.Params["cfgNil"].LeakageRate)
	assert.Equal(t, int64(100), manager.cfg.Params["cfgNil"].LenQueue)

	manager.cfg.Params = nil
	assert.NoError(t, manager.NewGroup("cfrParamsNil", params))

	groupParams = manager.GetParams("cfrParamsNil")
	assert.Equal(t, 6, len(manager.groupsMap))
	assert.Equal(t, 1, len(manager.cfg.Params))
	assert.Equal(t, int64(10), manager.cfg.Params["cfrParamsNil"].Capacity)
	assert.Equal(t, 1.0, manager.cfg.Params["cfrParamsNil"].LeakageRate)
	assert.Equal(t, int64(100), manager.cfg.Params["cfrParamsNil"].LenQueue)

}

func TestRecreateGroup(t *testing.T) {
	assert.NoError(t, initGlobalVars(t))
	tempDir := t.TempDir()
	configPathTempDir := filepath.Join(tempDir, "config.yaml")
	assert.Nil(t, os.WriteFile(configPathTempDir, []byte(configContent), 0644), "failed to write config file")
	manager, err := NewManager(tempDir)
	assert.NoError(t, err)
	require.NotNil(t, manager)

	groupParams := manager.GetParams("ip")
	assert.Equal(t, int64(10), groupParams.Capacity)
	assert.Equal(t, 16.666666666666668, groupParams.LeakageRate)
	assert.Equal(t, int64(10000), groupParams.LenQueue)

	params := cf.GroupParams{LenQueue: 100, Capacity: 5, LeakageRate: 1}
	manager.RecreateGroup("ip", params)

	groupParams = manager.GetParams("ip")
	assert.Equal(t, int64(5), groupParams.Capacity)
	assert.Equal(t, 1.0, groupParams.LeakageRate)
	assert.Equal(t, int64(100), groupParams.LenQueue)
}

func TestIsAllowed(t *testing.T) {
	require.NoError(t, initGlobalVars(t))

	name, id := "ip", "localhost"
	drops := int64(1)
	testErr := errors.New("Тестовая ошибка")

	specListMap[name].On("ToFind", id).Return(false, tp.IsBlack(false), nil).Once()
	groupsMap[name].On("IsAllowed", id, drops).Return(true).Once()
	assert.True(t, Manager.IsAllowed(name, id, drops))

	specListMap[name].On("ToFind", id).Return(false, tp.IsBlack(false), nil).Once()
	groupsMap[name].On("IsAllowed", id, drops).Return(false).Once()
	assert.False(t, Manager.IsAllowed(name, id, drops))

	specListMap[name].On("ToFind", id).Return(false, tp.IsBlack(false), testErr).Once()
	assert.False(t, Manager.IsAllowed(name, id, drops))

	specListMap[name].On("ToFind", id).Return(true, tp.IsBlack(false), nil).Once()
	assert.True(t, Manager.IsAllowed(name, id, drops))

	specListMap[name].On("ToFind", id).Return(true, tp.IsBlack(true), nil).Once()
	assert.False(t, Manager.IsAllowed(name, id, drops))
}

func TestAddToSpecList(t *testing.T) {
	require.NoError(t, initGlobalVars(t))

	name, id := "ip", "localhost"

	specListMock := specListMap[name]
	assert.False(t, Manager.AddToSpecList("err", id, tp.IsBlack(true)))

	specListMock.On("AddId", id, tp.IsBlack(true)).Return(true).Once()
	assert.True(t, Manager.AddToSpecList(name, id, tp.IsBlack(true)))

}

func TestDelFromSpecList(t *testing.T) {
	require.NoError(t, initGlobalVars(t))

	name, id := "ip", "localhost"

	specListMock := specListMap[name]
	specListMock.On("DelId", id).Return(true).Once()

	assert.False(t, Manager.DelFromSpecList("err", id))
	assert.True(t, Manager.DelFromSpecList(name, id))

	specListMock.On("DelId", id).Return(false).Once()
	assert.False(t, Manager.DelFromSpecList(name, id))
}

func TestIsRequestAllowed(t *testing.T) {
	require.NoError(t, initGlobalVars(t))

	// Все запросы разрешены
	for name := range Manager.groupsMap {
		specListMap[name].On("ToFind", "validID").Return(false, tp.IsBlack(false), nil).Once()
		groupsMap[name].On("IsAllowed", "validID", int64(1)).Return(true).Once()
	}

	params := tp.RequestParams{"login": "validID", "pass": "validID", "ip": "validID"}
	assert.True(t, Manager.IsRequestAllowed(params))

	// Один запрос запрещен
	for name := range Manager.groupsMap {
		if name == "ip" {
			specListMap["ip"].On("ToFind", "invalidID").Return(true, tp.IsBlack(true), nil).Once()
		} else {
			specListMap[name].On("ToFind", "validID").Return(false, tp.IsBlack(false), nil).Once()
			groupsMap[name].On("IsAllowed", "validID", int64(1)).Return(true).Once()
		}
	}
	params["ip"] = "invalidID"
	assert.False(t, Manager.IsRequestAllowed(params))

	// Ошибка при проверке списка
	specListMap["ip"].On("ToFind", "errorID").Return(false, tp.IsBlack(false), errors.New("ошибка")).Once()

	params = tp.RequestParams{"ip": "errorID"}
	assert.False(t, Manager.IsRequestAllowed(params))
}
