#!/usr/bin/python

import os
import random
import codecs
import platform
import shutil
import getopt
import sys

protoPackage = "protoMsg"

class Msg:
	name = ""
	fileType  = 0
	id = 0	
	def getName(self):
		pre = ""
		if self.fileType == 1:
			pre =  protoPackage + "."        
		return pre + self.name
	
	
def genMessageCmd(src_path,to_path):	
	maps = []
	route = {}
	out = open("out.go",'w')
	out.write("package common\n")
	out.write("import \"" + protoPackage + "\"\n")
	out.write("import \"zeus/msgdef\"\n")
	
	for filename in os.listdir(src_path):
		if not os.path.isdir(filename):
			fileType = 0
			if filename.endswith(".proto"):
				fileType = 1
			if filename.endswith(".binmsg"):
				fileType = 2
			if filename == "route":
				fileType = 3
			if fileType == 0:
				continue
			
			messageType = 1
			if filename.startswith("server"):
				messageType = 0
			
			fh = open(os.path.join(src_path , filename))
			for line in fh:
				line = line.replace(' ','').replace('{','').replace("\n","").replace("\t","")
				if line == "":
					continue
				if fileType == 3:
					t = line[:-1].split("=")
					if len(t) == 2:
						route[t[0]] = t[1]
				if fileType == 2 or fileType == 1:
					if line.startswith('message'):
						m =  Msg()
						m.name = line[7:]
						m.fileType = fileType
						m.id = 101
						if messageType == 1:
							m.id = m.id + 300
						maps.append(m)

	out.write("var ProtoMap = map[uint16]msgdef.IMsg{\n")	
	k = 1
	for value in maps:
		if value.name == "ProtoSync":
			out.write("10: new("+value.getName()+"),\n")
		else:
			id = str(value.id + k)
			if value.name in route:
				id = route[value.name] + id
			out.write(id+": new("+value.getName()+"),\n")
			value.id = int(id)
		k = k + 1
	out.write("}\n")	
	
	out.close()
	os.system("gofmt -l -w out.go")
	shutil.copyfile("out.go",to_path+"protoMap.go")
	os.remove("out.go")
	
	binProtoFile = open("proto.json",'w')
	binProtoFile.write(protoMapBinary(maps))
	binProtoFile.close()
	shutil.copyfile("proto.json", to_path+ "../../res/config/proto.json")
	os.remove("proto.json")

	genAllProtos(src_path, to_path)
	## gen game.proto.go
	# fh = open("game.proto", "rb")
	# temp1 = fh.read()
	# fh.close()
	# os.system("copy /y game.proto game.proto.bak")
	# temp1 = temp1.decode("gbk").encode("utf-8")
	# fh = open("game.proto", "wb")
	# fh.write(temp1)
	# fh.close()
	# os.system("protoc --gogofaster_out=" + to_path + "../../src/protoMsg/ game.proto"  )
	# os.system("copy /y game.proto.bak game.proto")
	# os.remove("game.proto.bak")
	
def genAllProtos(src_path, to_path):
	#backup & encode
	for filename in os.listdir(src_path):
		if not os.path.isdir(filename):
			if filename.endswith(".proto"):

				fh = open(filename, "rb")
				temp1 = fh.read()
				fh.close()
				
				os.system("copy /y " + filename + " " + filename + ".bak")
				
				temp1 = temp1.decode("gbk").encode("utf-8")
				fh = open(filename, "wb")
				fh.write(temp1)
				fh.close()
	
	#gen .proto.go
	for filename in os.listdir(src_path):
		if not os.path.isdir(filename):
			if filename.endswith(".proto"):
				os.system("protoc --gogofaster_out=" + to_path + "../../src/protoMsg/ " + filename)
				
	#restore
	for filename in os.listdir(src_path):
		if not os.path.isdir(filename):
			if filename.endswith(".proto.bak"):
				protoFileName = filename[:-4]				
				os.system("copy /y " + filename + " " + protoFileName)
				os.remove(filename)

def protoMapBinary(maps):
	data = "{"
	for value in maps:
		id = str(value.id)
		if value.name == "ProtoSync":
			id = "10"
		data += "\"" + id + "\":" + "\"" + value.name + "\","
	data = data[:-1]
	data += "}"
	return data
	
	
if __name__ == '__main__':
	realPath = os.path.dirname(os.path.realpath(__file__))
	src_path = realPath
	toPath = os.path.join(realPath ,"..\server\\src\common\\")
	genMessageCmd(src_path,toPath)
