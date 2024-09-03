package deep

import (
	"reflect"
)

func Copy[T any](src, dst *T) {
	srv := reflect.ValueOf(src)
	srv = srv.Elem()
	drv := reflect.ValueOf(dst)
	drv = drv.Elem()
	h := copyHandler{addrMap: map[uint64]reflect.Value{}}
	h.handle(srv, drv)
	clear(h.addrMap)
}

type copyHandler struct {
	addrMap map[uint64]reflect.Value
}

func (copyHandler) genKey(src reflect.Value) uint64 {
	switch src.Kind() {
	case reflect.Pointer, reflect.Chan, reflect.Map, reflect.UnsafePointer, reflect.Func, reflect.Slice:
		return uint64(src.Pointer())<<6 + uint64(src.Kind())<<1
	case reflect.Struct, reflect.Interface:
		return uint64(src.UnsafeAddr())<<6 + uint64(src.Kind())<<1 + 1
	default:
		return 0
	}
}

func (h copyHandler) handle(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		h.handleStruct(src, dst)
	case reflect.Interface:
		h.handleInterface(src, dst)
	case reflect.Pointer:
		h.handlePointer(src, dst)
	case reflect.Array:
		h.handleArray(src, dst)
	case reflect.Slice:
		h.handleSlice(src, dst)
	case reflect.Chan:
		dst.Set(reflect.MakeChan(src.Type(), src.Cap()))
	case reflect.Map:
		dst.Set(reflect.MakeMap(src.Type()))
		h.handleMap(src, dst)
	default:
		dst.Set(src)
	}
}

func (h copyHandler) handlePointer(src, dst reflect.Value) {
	src = src.Elem()
	addr := h.genKey(src)
	ndst, ok := h.addrMap[addr]
	if !ok {
		ndst = reflect.New(src.Type()).Elem()
		if addr > 0 {
			h.addrMap[addr] = ndst
		}
		h.handle(src, ndst)
	}
	dst.Set(ndst.Addr())
}

func (h copyHandler) handleStruct(src, dst reflect.Value) {
	srcAddr := h.genKey(src)
	if srcAddr > 0 {
		h.addrMap[srcAddr] = dst
	}
	for i := 0; i < src.NumField(); i++ {
		srcf := src.Field(i)
		addr := h.genKey(srcf)
		dstf, ok := h.addrMap[addr]
		if !ok {
			ndstf := dst.Field(i)
			if addr > 0 {
				h.addrMap[addr] = ndstf
			}
			h.handle(srcf, ndstf)
		} else {
			dst.Field(i).Set(dstf)
		}
	}
}

func (h copyHandler) handleInterface(src, dst reflect.Value) {
	src = src.Elem()
	h.handle(src, dst)
}

func (h copyHandler) handleArray(src, dst reflect.Value) {
	for i := 0; i < src.Len(); i++ {
		srcf := src.Index(i)
		addr := h.genKey(srcf)
		dstf, ok := h.addrMap[addr]
		if !ok {
			ndstf := dst.Index(i)
			if addr > 0 {
				h.addrMap[addr] = ndstf
			}
			h.handle(srcf, ndstf)
		} else {
			dst.Index(i).Set(dstf)
		}
	}
}

func (h copyHandler) handleSlice(src, dst reflect.Value) {
	srcAddr := h.genKey(src)
	if srcAddr > 0 {
		h.addrMap[srcAddr] = dst
	}
	dst.Grow(src.Len() - dst.Len())
	dst.SetLen(src.Len())
	for i := 0; i < src.Len(); i++ {
		srcf := src.Index(i)
		addr := h.genKey(srcf)
		dstf, ok := h.addrMap[addr]
		if !ok {
			ndstf := dst.Index(i)
			if addr > 0 {
				h.addrMap[addr] = ndstf
			}
			h.handle(srcf, ndstf)
		} else {
			dst.Index(i).Set(dstf)
		}
	}
}

func (h copyHandler) handleMap(src, dst reflect.Value) {
	srcAddr := h.genKey(src)
	if srcAddr > 0 {
		h.addrMap[srcAddr] = dst
	}
	iter := src.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		kAddr := h.genKey(k)
		kdst, ok := h.addrMap[kAddr]
		if !ok {
			kdst = reflect.New(k.Type()).Elem()
			if kAddr > 0 {
				h.addrMap[kAddr] = kdst
			}
			h.handle(k, kdst)
		}

		vAddr := h.genKey(v)
		vdst, ok := h.addrMap[vAddr]
		if !ok {
			vdst = reflect.New(v.Type()).Elem()
			if vAddr > 0 {
				h.addrMap[vAddr] = vdst
			}
			h.handle(v, vdst)
		}
		dst.SetMapIndex(kdst, vdst)
	}
}
