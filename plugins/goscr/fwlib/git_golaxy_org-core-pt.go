// Code generated by 'yaegi extract git.golaxy.org/core/pt'. DO NOT EDIT.

package fwlib

import (
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/pt"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/core/utils/option"
	"reflect"
)

func init() {
	Symbols["git.golaxy.org/core/pt/pt"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"CompWith":            reflect.ValueOf(pt.CompWith),
		"DefaultComponentLib": reflect.ValueOf(pt.DefaultComponentLib),
		"DefaultEntityLib":    reflect.ValueOf(pt.DefaultEntityLib),
		"EntityWith":          reflect.ValueOf(pt.EntityWith),
		"ErrPt":               reflect.ValueOf(&pt.ErrPt).Elem(),
		"For":                 reflect.ValueOf(pt.For),
		"NewComponentLib":     reflect.ValueOf(pt.NewComponentLib),
		"NewEntityLib":        reflect.ValueOf(pt.NewEntityLib),

		// type definitions
		"CompAtti":         reflect.ValueOf((*pt.CompAtti)(nil)),
		"ComponentDesc":    reflect.ValueOf((*pt.ComponentDesc)(nil)),
		"ComponentLib":     reflect.ValueOf((*pt.ComponentLib)(nil)),
		"ComponentPT":      reflect.ValueOf((*pt.ComponentPT)(nil)),
		"EntityAtti":       reflect.ValueOf((*pt.EntityAtti)(nil)),
		"EntityLib":        reflect.ValueOf((*pt.EntityLib)(nil)),
		"EntityPT":         reflect.ValueOf((*pt.EntityPT)(nil)),
		"EntityPTProvider": reflect.ValueOf((*pt.EntityPTProvider)(nil)),

		// interface wrapper definitions
		"_ComponentLib":     reflect.ValueOf((*_git_golaxy_org_core_pt_ComponentLib)(nil)),
		"_ComponentPT":      reflect.ValueOf((*_git_golaxy_org_core_pt_ComponentPT)(nil)),
		"_EntityLib":        reflect.ValueOf((*_git_golaxy_org_core_pt_EntityLib)(nil)),
		"_EntityPT":         reflect.ValueOf((*_git_golaxy_org_core_pt_EntityPT)(nil)),
		"_EntityPTProvider": reflect.ValueOf((*_git_golaxy_org_core_pt_EntityPTProvider)(nil)),
	}
}

// _git_golaxy_org_core_pt_ComponentLib is an interface wrapper for ComponentLib type
type _git_golaxy_org_core_pt_ComponentLib struct {
	IValue         interface{}
	WDeclare       func(comp any) pt.ComponentPT
	WGet           func(prototype string) (pt.ComponentPT, bool)
	WRange         func(fun generic.Func1[pt.ComponentPT, bool])
	WReversedRange func(fun generic.Func1[pt.ComponentPT, bool])
	WUndeclare     func(prototype string)
}

func (W _git_golaxy_org_core_pt_ComponentLib) Declare(comp any) pt.ComponentPT {
	return W.WDeclare(comp)
}
func (W _git_golaxy_org_core_pt_ComponentLib) Get(prototype string) (pt.ComponentPT, bool) {
	return W.WGet(prototype)
}
func (W _git_golaxy_org_core_pt_ComponentLib) Range(fun generic.Func1[pt.ComponentPT, bool]) {
	W.WRange(fun)
}
func (W _git_golaxy_org_core_pt_ComponentLib) ReversedRange(fun generic.Func1[pt.ComponentPT, bool]) {
	W.WReversedRange(fun)
}
func (W _git_golaxy_org_core_pt_ComponentLib) Undeclare(prototype string) {
	W.WUndeclare(prototype)
}

// _git_golaxy_org_core_pt_ComponentPT is an interface wrapper for ComponentPT type
type _git_golaxy_org_core_pt_ComponentPT struct {
	IValue      interface{}
	WConstruct  func() ec.Component
	WInstanceRT func() reflect.Type
	WPrototype  func() string
}

func (W _git_golaxy_org_core_pt_ComponentPT) Construct() ec.Component {
	return W.WConstruct()
}
func (W _git_golaxy_org_core_pt_ComponentPT) InstanceRT() reflect.Type {
	return W.WInstanceRT()
}
func (W _git_golaxy_org_core_pt_ComponentPT) Prototype() string {
	return W.WPrototype()
}

// _git_golaxy_org_core_pt_EntityLib is an interface wrapper for EntityLib type
type _git_golaxy_org_core_pt_EntityLib struct {
	IValue         interface{}
	WDeclare       func(prototype any, comps ...any) pt.EntityPT
	WGet           func(prototype string) (pt.EntityPT, bool)
	WGetEntityLib  func() pt.EntityLib
	WRange         func(fun generic.Func1[pt.EntityPT, bool])
	WReversedRange func(fun generic.Func1[pt.EntityPT, bool])
	WUndeclare     func(prototype string)
}

func (W _git_golaxy_org_core_pt_EntityLib) Declare(prototype any, comps ...any) pt.EntityPT {
	return W.WDeclare(prototype, comps...)
}
func (W _git_golaxy_org_core_pt_EntityLib) Get(prototype string) (pt.EntityPT, bool) {
	return W.WGet(prototype)
}
func (W _git_golaxy_org_core_pt_EntityLib) GetEntityLib() pt.EntityLib {
	return W.WGetEntityLib()
}
func (W _git_golaxy_org_core_pt_EntityLib) Range(fun generic.Func1[pt.EntityPT, bool]) {
	W.WRange(fun)
}
func (W _git_golaxy_org_core_pt_EntityLib) ReversedRange(fun generic.Func1[pt.EntityPT, bool]) {
	W.WReversedRange(fun)
}
func (W _git_golaxy_org_core_pt_EntityLib) Undeclare(prototype string) {
	W.WUndeclare(prototype)
}

// _git_golaxy_org_core_pt_EntityPT is an interface wrapper for EntityPT type
type _git_golaxy_org_core_pt_EntityPT struct {
	IValue              interface{}
	WAwakeOnFirstAccess func() *bool
	WComponent          func(idx int) pt.ComponentDesc
	WComponents         func() []pt.ComponentDesc
	WConstruct          func(settings ...option.Setting[ec.EntityOptions]) ec.Entity
	WCountComponents    func() int
	WInstanceRT         func() reflect.Type
	WPrototype          func() string
	WScope              func() *ec.Scope
}

func (W _git_golaxy_org_core_pt_EntityPT) AwakeOnFirstAccess() *bool {
	return W.WAwakeOnFirstAccess()
}
func (W _git_golaxy_org_core_pt_EntityPT) Component(idx int) pt.ComponentDesc {
	return W.WComponent(idx)
}
func (W _git_golaxy_org_core_pt_EntityPT) Components() []pt.ComponentDesc {
	return W.WComponents()
}
func (W _git_golaxy_org_core_pt_EntityPT) Construct(settings ...option.Setting[ec.EntityOptions]) ec.Entity {
	return W.WConstruct(settings...)
}
func (W _git_golaxy_org_core_pt_EntityPT) CountComponents() int {
	return W.WCountComponents()
}
func (W _git_golaxy_org_core_pt_EntityPT) InstanceRT() reflect.Type {
	return W.WInstanceRT()
}
func (W _git_golaxy_org_core_pt_EntityPT) Prototype() string {
	return W.WPrototype()
}
func (W _git_golaxy_org_core_pt_EntityPT) Scope() *ec.Scope {
	return W.WScope()
}

// _git_golaxy_org_core_pt_EntityPTProvider is an interface wrapper for EntityPTProvider type
type _git_golaxy_org_core_pt_EntityPTProvider struct {
	IValue        interface{}
	WGetEntityLib func() pt.EntityLib
}

func (W _git_golaxy_org_core_pt_EntityPTProvider) GetEntityLib() pt.EntityLib {
	return W.WGetEntityLib()
}