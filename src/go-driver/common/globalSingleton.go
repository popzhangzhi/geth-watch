package common

import (
	"sync"
)

/*
	单例实现全局配置
*/
type runtimeSingleton struct {
	//runtime 基础环境
	RuntimeEnv map[string]interface{}
	//runtime 业务
	RuntimeMain map[string]interface{}
	//地址秘钥字典
	addresses *addresses
}

//地址对
type addresses struct {
	//系统地址
	systems [][]byte
	//全地址
	allAddresses [][]byte
	//全地址索引的秘钥
	allKey map[string][]byte
}

var instance *runtimeSingleton
var once sync.Once

func GetInstance() *runtimeSingleton {
	once.Do(func() {

		address := &addresses{}
		instance = &runtimeSingleton{
			make(map[string]interface{}),
			make(map[string]interface{}),
			address,
		}
	})
	return instance
}

/*
	设置 整个runtimeEnv
*/
func (R *runtimeSingleton) SetEnvMap(config map[string]interface{}) {
	R.RuntimeEnv = config
}

/*
	设置 runtimeEnv
*/
func (R *runtimeSingleton) SetEnv(name string, value interface{}) {

	R.RuntimeEnv[name] = value
}

/*
	获取 runtimeEnv
*/
func (R *runtimeSingleton) GetEnv(name string) interface{} {
	return R.RuntimeEnv[name]

}

func (R *runtimeSingleton) SetAddresses(systems [][]byte, allAddresses [][]byte, key map[string][]byte) {

	R.addresses.systems = systems
	R.addresses.allAddresses = allAddresses
	R.addresses.allKey = key
}

func (R *runtimeSingleton) GetKey(address string) []byte {

	return R.addresses.allKey[address]
}

func (R *runtimeSingleton) GetSystems() [][]byte {
	return R.addresses.systems
}
func (R *runtimeSingleton) GetAllAddresses() [][]byte {
	return R.addresses.allAddresses
}
