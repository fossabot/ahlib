package xdi

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"reflect"
	"sync"
)

// DiContainer log function.
// Yellow color used for type (field and parent), red color used for name (field)
//
// kind: Name | Impl | Type | Inject
// parentName: _ if not Inject else parentType
type LogFunc func(kind, parentType, fieldName, fieldType string)

// Global service name type (prevent package import cycle).
type ServiceName string

func (s *ServiceName) String() string {
	return string(*s)
}

type DiContainer struct {
	provByType map[reflect.Type]interface{}
	provByName map[ServiceName]interface{}
	muByType   sync.Mutex
	muByName   sync.Mutex

	provideLog bool
	injectLog  bool

	// Log function, just like:
	//	[XDI] Name:    a (*xdi.ServiceA)                  -> RED{a} YELLOW{*xdi.ServiceA}
	//	[XDI] Impl:    _ (xdi.IServiceA)                  -> RED{_} YELLOW{xdi.IServiceA}
	//	[XDI] Type:    _ (*xdi.ServiceB)                  -> RED{_} YELLOW{*xdi.ServiceB}
	//	[XDI] Inject:  (*xdi.ServiceB).SA (xdi.IServiceA) -> YELLOW{*xdi.ServiceB} RED{SA} YELLOW{xdi.IServiceA}
	logFunc LogFunc
}

func NewDiContainer() *DiContainer {
	return &DiContainer{
		provByType: make(map[reflect.Type]interface{}),
		provByName: make(map[ServiceName]interface{}),
		provideLog: true,
		injectLog:  true,
		logFunc:    DefaultLogFunc(),
	}
}

func (d *DiContainer) SetLogMode(provideLog bool, injectLog bool) {
	_di.provideLog = provideLog
	_di.injectLog = injectLog
}

func (d *DiContainer) SetLogFunc(logFunc LogFunc) {
	_di.logFunc = logFunc
}

func DefaultLogFunc() LogFunc {
	return func(kind, parentType, fieldName, fieldType string) {
		xcolor.ForceColor()
		kind += ":"
		if parentType != "" {
			parentType = fmt.Sprintf("(%s).", xcolor.Yellow.Sprint(parentType))
		}
		fieldName = xcolor.Red.Sprint(fieldName)
		fieldType = xcolor.Yellow.Sprint(fieldType)
		fmt.Printf("[XDI] %-8s %s%s (%s)\n", kind, parentType, fieldName, fieldType)
	}
}

func (d *DiContainer) ProvideType(service interface{}) {
	if service == nil {
		panic("could not provide nil service")
	}
	t := reflect.TypeOf(service)

	d.muByType.Lock()
	d.provByType[t] = service
	d.muByType.Unlock()
	if d.provideLog {
		d.logFunc("Type", "", "_", t.String())
	}
}

func (d *DiContainer) ProvideImpl(itfPtr interface{}, impl interface{}) {
	it := reflect.TypeOf(itfPtr)
	if reflect.TypeOf(it).Kind() != reflect.Ptr {
		panic("first parameter of ProvideImpl must be pointer of interface")
	}
	it = it.Elem()

	st := reflect.TypeOf(impl)
	if !st.Implements(it) {
		panic(fmt.Sprintf("could not implement type %s by %s", it.String(), st.String()))
	}

	d.muByType.Lock()
	d.provByType[it] = impl
	d.muByType.Unlock()
	if d.provideLog {
		d.logFunc("Impl", "", "_", it.String())
	}
}

func (d *DiContainer) ProvideName(name ServiceName, service interface{}) {
	if name == "~" {
		panic("could not provide service using '~' name")
	}

	d.muByName.Lock()
	d.provByName[name] = service
	d.muByName.Unlock()
	if d.provideLog {
		d.logFunc("Name", "", string(name), reflect.TypeOf(service).String())
	}
}

func (d *DiContainer) GetByType(srvType interface{}) (service interface{}, exist bool) {
	if srvType == nil {
		panic("could not get nil type service")
	}
	service, exist = d.provByType[reflect.TypeOf(srvType)]
	return
}

func (d *DiContainer) GetByName(name ServiceName) (service interface{}, exist bool) {
	service, exist = d.provByName[name]
	return
}

func (d *DiContainer) GetByTypeForce(srvType interface{}) interface{} {
	service, exist := d.GetByType(srvType)
	if !exist {
		panic(fmt.Sprintf("service with type %s is not found", reflect.TypeOf(srvType).String()))
	}
	return service
}

func (d *DiContainer) GetByNameForce(name ServiceName) interface{} {
	service, exist := d.GetByName(name)
	if !exist {
		panic(fmt.Sprintf("service with name %s is not found", name))
	}
	return service
}

func (d *DiContainer) Inject(ctrl interface{}) (allInjected bool) {
	return d.inject(ctrl, false)
}

func (d *DiContainer) MustInject(ctrl interface{}) {
	d.inject(ctrl, true)
}

// Using tips:
//
//	`di:""`       // -> ignore
//	`di:"-"`      // -> ignore
//	`di:"~"`      // -> auto inject
//	`di:"name"`   // -> inject by name
func (d *DiContainer) inject(ctrl interface{}, force bool) bool {
	ctrlType := xreflect.ElemType(ctrl)
	ctrlValue := xreflect.ElemValue(ctrl)
	if ctrlType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("object for injection should be struct, have %s with %s", ctrlType.Kind().String(), ctrlType.String()))
	}
	allInjected := true

	for fieldIdx := 0; fieldIdx < ctrlType.NumField(); fieldIdx++ {
		// check
		field := ctrlType.Field(fieldIdx)
		diTag := field.Tag.Get("di")
		if diTag == "-" || diTag == "" {
			continue
		}

		// find
		var service interface{}
		var exist bool
		if diTag == "~" {
			service, exist = d.provByType[field.Type]
		} else {
			service, exist = d.provByName[ServiceName(diTag)]
		}

		// exist
		if !exist {
			allInjected = false
			if force {
				panic("there are some fields could not be injected")
			}
			continue
		}

		// inject
		fieldType := ctrlType.Field(fieldIdx)
		fieldValue := ctrlValue.Field(fieldIdx)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			fieldValue.Set(reflect.ValueOf(service))
			if d.injectLog {
				d.logFunc("Inject", reflect.TypeOf(ctrl).String(), fieldType.Name, fieldType.Type.String())
			}
		}
	}

	return allInjected
}

// A DiContainer that used for global
var _di = NewDiContainer()

// A global SetLogMode function for DiContainer
// only for global DiContainer, not work for NewDiContainer
func SetLogMode(provideLogMode bool, injectLogMode bool) {
	_di.provideLog = provideLogMode
	_di.injectLog = injectLogMode
}

// A global SetLogFunc function for DiContainer
// only for global DiContainer, not work for NewDiContainer
func SetLogFunc(logFunc LogFunc) {
	_di.logFunc = logFunc
}

func ProvideType(service interface{}) {
	_di.ProvideType(service)
}

func ProvideImpl(interfacePtr interface{}, impl interface{}) {
	_di.ProvideImpl(interfacePtr, impl)
}

func ProvideName(name ServiceName, service interface{}) {
	_di.ProvideName(name, service)
}

func GetByType(srvType interface{}) (service interface{}, exist bool) {
	return _di.GetByType(srvType)
}

func GetByName(name ServiceName) (service interface{}, exist bool) {
	return _di.GetByName(name)
}

func GetByTypeForce(srvType interface{}) interface{} {
	return _di.GetByTypeForce(srvType)
}

func GetByNameForce(name ServiceName) interface{} {
	return _di.GetByNameForce(name)
}

func Inject(ctrl interface{}) (allInjected bool) {
	return _di.Inject(ctrl)
}

func MustInject(ctrl interface{}) {
	_di.MustInject(ctrl)
}
