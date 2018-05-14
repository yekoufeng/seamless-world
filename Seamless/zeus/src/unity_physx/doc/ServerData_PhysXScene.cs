using UnityEngine;
using UnityEditor;
using System.IO;
using System.Collections.Generic;

partial class ServerData {

    [MenuItem("ServerData/Export PhysXScene")]
    public static void exportPhysXScene() {

        var path = EditorUtility.SaveFilePanel("Export PhysXScene", Application.dataPath + "/..", "pxscene", string.Empty);
        if (path == null || path.Length == 0) {
            return;
        }

        System.Collections.Generic.List<GameObject> terrainTrees = new System.Collections.Generic.List<GameObject>();
        foreach (var terrain in Terrain.activeTerrains) {
            var tt = new GameObject("__Trees");
            terrainTrees.Add(tt);
            tt.transform.SetParent(tt.transform, false);
            var td = terrain.terrainData;
            foreach (var ti in td.treeInstances) {
                var tp = td.treePrototypes[ti.prototypeIndex];
                var go = Object.Instantiate<GameObject>(tp.prefab, tt.transform, false);
                go.transform.localPosition = Vector3.Scale(ti.position, td.size);
                go.transform.localRotation = Quaternion.Euler(0, ti.rotation, 0);
                go.transform.localScale = new Vector3(ti.widthScale, ti.heightScale, ti.widthScale);
            }
        }

        HashSet<Collider> m_excludes = new HashSet<Collider>();
        foreach (var dv in Object.FindObjectsOfType<DoorView>()) {
            foreach (var col in dv.GetComponentsInChildren<Collider>()) {
                m_excludes.Add(col);
            }
        }
        

        PxMeshDictionary meshes = new PxMeshDictionary();
        List<PxSceneObject> objs = new System.Collections.Generic.List<PxSceneObject>();

        foreach (var box in Object.FindObjectsOfType<BoxCollider>()) {
            if (!m_excludes.Contains(box)) {
                objs.Add(new PxBoxCollider(box));
            }
        }

        foreach (var capsule in Object.FindObjectsOfType<CapsuleCollider>()) {
            if (!m_excludes.Contains(capsule)) {
                objs.Add(new PxCapsuleCollider(capsule));
            }
        }

        foreach (var mesh in Object.FindObjectsOfType<MeshCollider>()) {
            if (!mesh.convex && mesh.sharedMesh != null) {
                if (!m_excludes.Contains(mesh)) {
                    objs.Add(new PxMeshCollider(mesh, meshes));
                }
            }
        }

        foreach (var terrain in Object.FindObjectsOfType<TerrainCollider>()) {
            if (terrain.terrainData != null) {
                if (!m_excludes.Contains(terrain)) {
                    objs.Add(new PxTerrainCollider(terrain));
                }
            }
        }

        foreach (var tt in terrainTrees) {
            Object.DestroyImmediate(tt);
        }
        terrainTrees.Clear();

        using (var file = new FileStream(path, FileMode.Create)) {
            var bw = new BinaryWriter(file);
            bw.Write(new byte[] { (byte)'P', (byte)'X', (byte)'S', 0 });
            var meshArray = meshes.toArray();
            bw.Write(meshArray.Length);
            for (int i = 0; i < meshArray.Length; ++i) {
                meshArray[i].save(bw);
            }
            bw.Write(objs.Count);
            for (int i = 0; i < objs.Count; ++i) {
                objs[i].save(bw);
            }
        }

        EditorUtility.DisplayDialog("Export PhysXScene", "Success", "Ok");

        //var root = new GameObject("PhysXScene");

        //foreach (var mesh in meshes.toArray()) {
        //    mesh.dump(root.transform);
        //}

        //foreach (var obj in objs) {
        //    obj.dump(root.transform);
        //}
    }


    enum PxObjectType : ushort {
        kMesh = 1,
        kBoxCollider = 2,
        kCapsuleCollider = 3,
        kMeshCollider = 4,
        kTerrainCollider = 5,
    }

    abstract class PxObject {
        public abstract PxObjectType type { get; }
        public virtual void save(System.IO.BinaryWriter bw) {
            bw.Write((ushort)type);
        }

        public virtual void dump(Transform root) { }
    }

    abstract class PxSceneObject : PxObject {
        public Vector3 position;
        public Quaternion rotation;
        public byte layer;

        protected void _setPositionAndRotation(Transform transform) {
            position = transform.position;
            rotation = transform.rotation;
            layer = (byte)transform.gameObject.layer;
        }

        public override void save(BinaryWriter bw) {
            base.save(bw);
            bw.Write(position.x);
            bw.Write(position.y);
            bw.Write(position.z);
            bw.Write(rotation.x);
            bw.Write(rotation.y);
            bw.Write(rotation.z);
            bw.Write(rotation.w);
            bw.Write(layer);
        }
    }

    class PxBoxCollider : PxSceneObject {
        public static readonly PxObjectType clsType = PxObjectType.kBoxCollider;
        public override PxObjectType type { get { return clsType; } }
        public Vector3 halfExtents;

        public PxBoxCollider() { }

        public PxBoxCollider(BoxCollider source) {
            position = source.transform.TransformPoint(source.center);
            rotation = source.transform.rotation;
            layer = (byte)source.transform.gameObject.layer;
            halfExtents = Vector3.Scale(source.size, source.transform.lossyScale) * 0.5f;
        }

        public override void save(BinaryWriter bw) {
            base.save(bw);
            bw.Write(halfExtents.x);
            bw.Write(halfExtents.y);
            bw.Write(halfExtents.z);
        }

        public override void dump(Transform root) {
            var go = new GameObject("box");
            go.transform.SetPositionAndRotation(position, rotation);
            go.layer = layer;
            var box = go.AddComponent<BoxCollider>();
            box.size = halfExtents * 2;

            go.transform.SetParent(root, true);
        }
    }

    class PxCapsuleCollider : PxSceneObject {
        public static readonly PxObjectType clsType = PxObjectType.kCapsuleCollider;
        public override PxObjectType type { get { return clsType; } }
        public float radius;
        public float halfHeight;

        public PxCapsuleCollider() { }
        public PxCapsuleCollider(CapsuleCollider source) {
            position = source.transform.TransformPoint(source.center);
            rotation = source.transform.rotation;
            layer = (byte)source.transform.gameObject.layer;
            var scale = source.transform.lossyScale;
            var height = 0f;
            switch (source.direction) {
                case 0: // x-axis
                    radius = source.radius * Mathf.Max(scale.y, scale.z);
                    height = source.height * scale.x;
                    break;
                case 1: // y-axis
                    rotation *= Quaternion.Euler(0, 0, -90);
                    radius = source.radius * Mathf.Max(scale.x, scale.z);
                    height = source.height * scale.y;

                    break;
                case 2: // z-axis
                    rotation *= Quaternion.Euler(0, 90, 0);
                    radius = source.radius * Mathf.Max(scale.x, scale.y);
                    height = source.height * scale.z;
                    break;
            }
            halfHeight = Mathf.Max((height - radius * 2) * 0.5f, 0);
        }

        public override void save(BinaryWriter bw) {
            base.save(bw);
            bw.Write(radius);
            bw.Write(halfHeight);
        }

        public override void dump(Transform root) {
            var go = new GameObject("capsule");
            go.transform.SetPositionAndRotation(position, rotation);
            go.layer = layer;
            var capsule = go.AddComponent<CapsuleCollider>();
            capsule.height = (radius + halfHeight) * 2;
            capsule.radius = radius;

            go.transform.SetParent(root, true);
        }
    }

    class PxMesh : PxObject {
        public static readonly PxObjectType clsType = PxObjectType.kMesh;
        public override PxObjectType type { get { return clsType; } }
        public Vector3[] vertices;
        public ushort[] indices;

        public PxMesh(int referenceIndex) { this.referenceIndex = referenceIndex; }

        public int referenceIndex { get; private set; }
        public Mesh mesh { get; private set; }

        public override void save(BinaryWriter bw) {
            base.save(bw);
            bw.Write(vertices.Length);
            for (int i = 0; i < vertices.Length; ++i) {
                var v = vertices[i];
                bw.Write(v.x);
                bw.Write(v.y);
                bw.Write(v.z);
            }
            bw.Write(indices.Length);
            for (int i = 0; i < indices.Length; ++i) {
                bw.Write(indices[i]);
            }
        }

        public override void dump(Transform root) {
            mesh = new Mesh();
            mesh.vertices = vertices;
            var indices = new int[this.indices.Length];
            for (int i = 0; i < indices.Length; ++i) {
                indices[i] = this.indices[i];
            }
            mesh.SetIndices(indices, MeshTopology.Triangles, 0);
        }
    }

    class PxMeshDictionary {
        
        public PxMesh buildMesh(Mesh mesh) {
            PxMesh ret;
            if (!m_index.TryGetValue(mesh, out ret)) {
                ret = new PxMesh(m_meshes.Count);
                ret.vertices = mesh.vertices;
                uint indexCount = 0;
                for (int i = 0; i < mesh.subMeshCount; ++i) {
                    indexCount += mesh.GetIndexCount(i);
                }
                ret.indices = new ushort[indexCount];
                int index = 0;
                for (int i = 0; i < mesh.subMeshCount; ++i) {
                    var indices = mesh.GetIndices(i);
                    for (int j = 0; j < indices.Length; ++j) {
                        ret.indices[index++] = (ushort)indices[j];
                    }
                }
                m_meshes.Add(ret);
                m_index.Add(mesh, ret);
            }
            return ret;
        }

        public PxMesh[] toArray() { return m_meshes.ToArray(); }

        System.Collections.Generic.Dictionary<Mesh, PxMesh> m_index = new System.Collections.Generic.Dictionary<Mesh, PxMesh>();
        System.Collections.Generic.List<PxMesh> m_meshes = new System.Collections.Generic.List<PxMesh>();
    }

    class PxMeshCollider : PxSceneObject {
        public static readonly PxObjectType clsType = PxObjectType.kMeshCollider;
        public override PxObjectType type { get { return clsType; } }

        public Vector3 scale;
        public PxMesh mesh;

        public PxMeshCollider() { }
        public PxMeshCollider(MeshCollider source, PxMeshDictionary meshDictionary) {
            _setPositionAndRotation(source.transform);
            scale = source.transform.lossyScale;
            mesh = meshDictionary.buildMesh(source.sharedMesh);
        }

        public override void save(BinaryWriter bw) {
            base.save(bw);
            bw.Write(scale.x);
            bw.Write(scale.y);
            bw.Write(scale.z);
            bw.Write(mesh.referenceIndex);
        }

        public override void dump(Transform root) {
            var go = new GameObject("mesh");
            go.transform.SetPositionAndRotation(position, rotation);
            go.transform.localScale = scale;
            go.layer = layer;
            var mesh = go.AddComponent<MeshCollider>();
            mesh.sharedMesh = this.mesh.mesh;

            go.transform.SetParent(root, true);
        }
    }

    class PxTerrainCollider : PxSceneObject {
        public static readonly PxObjectType clsType = PxObjectType.kTerrainCollider;
        public override PxObjectType type { get { return clsType; } }

        public Vector3 size;
        public float[,] heightmap;

        public PxTerrainCollider() { }
        public PxTerrainCollider(TerrainCollider source) {
            _setPositionAndRotation(source.transform);
            size = Vector3.Scale(source.terrainData.size, source.transform.lossyScale);
            heightmap = source.terrainData.GetHeights(0, 0, source.terrainData.heightmapWidth, source.terrainData.heightmapHeight);
        }

        public override void save(BinaryWriter bw) {
            base.save(bw);
            bw.Write(size.x);
            bw.Write(size.y);
            bw.Write(size.z);
            var d = heightmap.GetLength(0);
            bw.Write(d);
            for (int i = 0; i < d; ++i) {
                for (int j = 0; j < d; ++j) {
                    bw.Write(heightmap[i, j]);
                }
            }
        }
    }
}
