package main

import (
	"errors"
	"io/ioutil"
	"zeus/common"
	"zeus/nav"
)

func (ms *Maps) loadNavMesh(path string) (*nav.Mesh, *nav.MeshPathFinder, error) {

	mesh := nav.NewMesh(path + "navmesh.bin")
	if mesh == nil {
		return nil, nil, errors.New("load nav mesh fail")
	}

	pathFinder := nav.NewMeshPathFinder(mesh)
	return mesh, pathFinder, nil
}

func (ms *Maps) loadRange(path string) (*MapRanges, error) {

	data, err := ioutil.ReadFile(path + "spawn_probes.probes")
	if err != nil {
		return nil, err
	}

	rs := newMapRanges()

	br := common.NewByteStream(data)

	b1, _ := br.ReadByte()
	b2, _ := br.ReadByte()
	b3, _ := br.ReadByte()
	b4, _ := br.ReadByte()

	if !(b1 == 'S' && b2 == 'P' && b3 == 'T' && b4 == 0) {
		return nil, errors.New("wrong ranges file header")
	}

	num, _ := br.ReadUInt32()

	for i := 0; i < int(num); i++ {
		typ, _ := br.ReadUInt16()
		id, _ := br.ReadInt32()
		centerX, _ := br.ReadFloat32()
		centerY, _ := br.ReadFloat32()
		centerZ, _ := br.ReadFloat32()
		radius, _ := br.ReadFloat32()

		r := newMapRange()
		r.Typ = int(typ)
		r.ID = int(id)
		r.CenterPos.X = centerX
		r.CenterPos.Y = centerY
		r.CenterPos.Z = centerZ
		r.Radius = radius

		rs.addRange(r)
	}

	return rs, nil
}

func (ms *Maps) loadHeight(path string) (*MapHeightMap, error) {

	data, err := ioutil.ReadFile(path + "height_map.map")
	if err != nil {
		return nil, err
	}

	br := common.NewByteStream(data)

	b1, _ := br.ReadByte()
	b2, _ := br.ReadByte()
	b3, _ := br.ReadByte()
	b4, _ := br.ReadByte()

	if !(b1 == 'T' && b2 == 'H' && b3 == 'M' && b4 == 0) {
		return nil, errors.New("wrong height map file header")
	}
	ox, _ := br.ReadFloat32()
	oz, _ := br.ReadFloat32()

	w, _ := br.ReadFloat32()
	h, _ := br.ReadFloat32()

	r, _ := br.ReadUInt32()

	mh := newMapHeightMap(w, h, ox, oz, r)

	for i := 0; i < int(r); i++ {

		for j := 0; j < int(r); j++ {
			h, _ := br.ReadFloat32()
			mh.Data[i*int(r)+j] = h
		}
	}

	return mh, nil
}
