package manager

import (
	"fmt"
	"sync"

	cf "bruteforce/config"
	sr "bruteforce/internal"
	bk "bruteforce/internal/bucket"
	gr "bruteforce/internal/group"
	sl "bruteforce/internal/listsystem"
	tp "bruteforce/internal/types"
	"github.com/fat0troll/durufmt"
)

type groupMap map[string]GroupItem

//go:generate mockery --name GroupItem
type GroupItem interface {
	GetCapacity() int64
	GetLeakageRate() float64
	GetLenQueue() int64
	IsAllowed(id string, drops int64) (isAllowed bool)
	RemoveDrops(id string)
	GetAllGroupData() *tp.AllGroupData
}

//go:generate mockery --name SpecListsManager
type SpecListsManager interface {
	GetName() string
	AddID(id string, isBlack tp.IsBlack) (ok bool)
	DelID(id string) (ok bool)
	ToFind(id string) (isFound bool, isBlack tp.IsBlack, err error)
	Reset()
	GetIDsMap() map[string]tp.IsBlack
}

type GroupManager struct {
	lock        sync.Mutex                  // Для потокобезопастности
	cfg         *cf.Config                  // Параметры групп из конфига
	groupsMap   groupMap                    // Карта групп
	specListMap map[string]SpecListsManager // Карта СпецСписков
}

func NewManager(path string) (*GroupManager, error) {
	cfg, err := cf.LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке конфига: path = %s, err - %w", path, err)
	}
	manager := GroupManager{
		cfg:         &cfg,
		groupsMap:   groupMap{},
		specListMap: map[string]SpecListsManager{},
	}
	for name, params := range manager.cfg.Params {
		manager.groupsMap[name] = gr.NewBucketGroup(
			params.LeakageRate,
			params.Capacity,
			params.LenQueue,
			bk.NewQueyeBuckets(int(params.LenQueue)),
			true,
			bk.NewBucket,
			bk.NewQueyeBuckets,
		)
		if name == "ip" {
			manager.specListMap[name] = sl.NewSpecList(name, 100000, sl.SearchEngineIP, sl.SubnetNameVerifier)
		} else {
			manager.specListMap[name] = sl.NewSpecList(name, 100000, sl.SearchEngineCommon, sl.CommonNameVerifier)
		}
	}
	return &manager, nil
}

func (g *GroupManager) NewGroup(name string, params cf.GroupParams) error {
	if g.cfg == nil {
		g.cfg = &cf.Config{Params: map[string]cf.GroupParams{name: params}}
		return g.createBucketGroup(name, params) // Создаем группу и выходим
	}
	if g.cfg.Params == nil {
		g.cfg.Params = map[string]cf.GroupParams{name: params}
		return g.createBucketGroup(name, params) // Создаем группу и выходим
	}
	if _, ok := g.cfg.Params[name]; ok {
		return fmt.Errorf("такая группа уже существует: name = %s", name)
	}
	g.cfg.Params[name] = params
	return g.createBucketGroup(name, params)
}

func (g *GroupManager) createBucketGroup(name string, params cf.GroupParams) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.groupsMap[name] = gr.NewBucketGroup(
		params.LeakageRate,
		params.Capacity,
		params.LenQueue,
		bk.NewQueyeBuckets(int(params.LenQueue)),
		true,
		bk.NewBucket,
		bk.NewQueyeBuckets,
	)
	if name == "ip" {
		g.specListMap[name] = sl.NewSpecList(name, 100000, sl.SearchEngineIP, sl.SubnetNameVerifier)
	} else {
		g.specListMap[name] = sl.NewSpecList(name, 100000, sl.SearchEngineCommon, sl.CommonNameVerifier)
	}
	return nil
}

func (g *GroupManager) RecreateGroup(name string, params cf.GroupParams) { // Через запрос
	if _, ok := g.cfg.Params[name]; ok {
		gr.DelBucketGroup(g.groupsMap[name].(*gr.BucketGroup))
	}
	g.lock.Lock()
	defer g.lock.Unlock()

	g.groupsMap[name] = gr.NewBucketGroup(
		params.LeakageRate,
		params.Capacity,
		params.LenQueue,
		bk.NewQueyeBuckets(int(params.LenQueue)),
		true,
		bk.NewBucket,
		bk.NewQueyeBuckets,
	)
	if name == "ip" {
		g.specListMap[name] = sl.NewSpecList(name, 100000, sl.SearchEngineIP, sl.SubnetNameVerifier)
	} else {
		g.specListMap[name] = sl.NewSpecList(name, 100000, sl.SearchEngineCommon, sl.CommonNameVerifier)
	}
}

func (g *GroupManager) GetParams(name string) (params cf.GroupParams) {
	params.LeakageRate = g.groupsMap[name].GetLeakageRate()
	params.Capacity = g.groupsMap[name].GetCapacity()
	params.LenQueue = g.groupsMap[name].GetLenQueue()
	return
}

func (g *GroupManager) AddToSpecList(name, id string, isBlack tp.IsBlack) (ok bool) {
	if sList, ok := g.specListMap[name]; ok {
		g.lock.Lock()
		defer g.lock.Unlock()
		return sList.AddID(id, isBlack)
	}
	return false
}

func (g *GroupManager) DelFromSpecList(name, id string) (ok bool) {
	if sList, ok := g.specListMap[name]; ok {
		g.lock.Lock()
		defer g.lock.Unlock()
		return sList.DelID(id)
	}
	return false
}

func (g *GroupManager) IsAllowed(nameGroup, idBucket string, drops int64) (isAllowed bool) {
	isFound, isBlack, err := g.specListMap[nameGroup].ToFind(idBucket)
	if err != nil {
		return false
	}
	if isFound {
		return !bool(isBlack)
	}
	if group, ok := g.groupsMap[nameGroup]; ok {
		return group.IsAllowed(idBucket, drops)
	}
	return false
}

func (g *GroupManager) IsRequestAllowed(params tp.RequestParams) (isAllowed bool) {
	var wg sync.WaitGroup
	results := make(chan bool, len(params))

	for name, id := range params {
		wg.Add(1)
		go func(name, id string, wg *sync.WaitGroup, results chan<- bool) {
			defer wg.Done()
			results <- g.IsAllowed(name, id, int64(1))
		}(name, id, &wg, results)
	}

	wg.Wait()
	close(results)

	finalResult := true
	for result := range results {
		finalResult = finalResult && result
	}
	return finalResult
}

func (g *GroupManager) ResetBacket(name, id string) (err error) {
	if group, ok := g.groupsMap[name]; ok {
		g.lock.Lock()
		defer g.lock.Unlock()
		group.RemoveDrops(id)
		return nil
	}
	return fmt.Errorf("группа с наименованием %v не найдена", name)
}

func (g *GroupManager) GetAllGroupsData() tp.AllGroupsData {
	groups := tp.AllGroupsData{}
	groups.Elapsed = fmt.Sprintf("Текущий сдвиг 'elapsed' = %v", durufmt.Parse(sr.GetElapsed()).String())
	groups.LenGroupsMap = len(g.groupsMap)
	groups.GroupsParams = make(map[string]tp.AllGroupData, groups.LenGroupsMap)
	for name, group := range g.groupsMap {
		if params := group.GetAllGroupData(); params != nil {
			groups.GroupsParams[name] = *params
		}
	}
	groups.LenSList = len(g.specListMap)
	groups.SpecList = map[string]map[string]tp.IsBlack{}
	for name, sListManager := range g.specListMap {
		groups.SpecList[name] = sListManager.GetIDsMap()
	}
	return groups
}
