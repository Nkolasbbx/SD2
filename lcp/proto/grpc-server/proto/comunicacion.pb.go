// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: comunicacion.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Entrenador struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Nombre        string                 `protobuf:"bytes,2,opt,name=nombre,proto3" json:"nombre,omitempty"`
	Region        string                 `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	Ranking       int32                  `protobuf:"varint,4,opt,name=ranking,proto3" json:"ranking,omitempty"`
	Estado        string                 `protobuf:"bytes,5,opt,name=estado,proto3" json:"estado,omitempty"`
	Suspension    int32                  `protobuf:"varint,6,opt,name=suspension,proto3" json:"suspension,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Entrenador) Reset() {
	*x = Entrenador{}
	mi := &file_comunicacion_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Entrenador) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entrenador) ProtoMessage() {}

func (x *Entrenador) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entrenador.ProtoReflect.Descriptor instead.
func (*Entrenador) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{0}
}

func (x *Entrenador) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Entrenador) GetNombre() string {
	if x != nil {
		return x.Nombre
	}
	return ""
}

func (x *Entrenador) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Entrenador) GetRanking() int32 {
	if x != nil {
		return x.Ranking
	}
	return 0
}

func (x *Entrenador) GetEstado() string {
	if x != nil {
		return x.Estado
	}
	return ""
}

func (x *Entrenador) GetSuspension() int32 {
	if x != nil {
		return x.Suspension
	}
	return 0
}

type Torneo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IdTorneo      int64                  `protobuf:"varint,1,opt,name=idTorneo,proto3" json:"idTorneo,omitempty"`
	Region        string                 `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	Jugadores     []*Entrenador          `protobuf:"bytes,3,rep,name=jugadores,proto3" json:"jugadores,omitempty"`
	Estado        string                 `protobuf:"bytes,5,opt,name=estado,proto3" json:"estado,omitempty"`
	Inscritos     int32                  `protobuf:"varint,6,opt,name=inscritos,proto3" json:"inscritos,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Torneo) Reset() {
	*x = Torneo{}
	mi := &file_comunicacion_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Torneo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Torneo) ProtoMessage() {}

func (x *Torneo) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Torneo.ProtoReflect.Descriptor instead.
func (*Torneo) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{1}
}

func (x *Torneo) GetIdTorneo() int64 {
	if x != nil {
		return x.IdTorneo
	}
	return 0
}

func (x *Torneo) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Torneo) GetJugadores() []*Entrenador {
	if x != nil {
		return x.Jugadores
	}
	return nil
}

func (x *Torneo) GetEstado() string {
	if x != nil {
		return x.Estado
	}
	return ""
}

func (x *Torneo) GetInscritos() int32 {
	if x != nil {
		return x.Inscritos
	}
	return 0
}

type ListaTorneos struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Torneos       []*Torneo              `protobuf:"bytes,1,rep,name=torneos,proto3" json:"torneos,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListaTorneos) Reset() {
	*x = ListaTorneos{}
	mi := &file_comunicacion_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListaTorneos) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListaTorneos) ProtoMessage() {}

func (x *ListaTorneos) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListaTorneos.ProtoReflect.Descriptor instead.
func (*ListaTorneos) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{2}
}

func (x *ListaTorneos) GetTorneos() []*Torneo {
	if x != nil {
		return x.Torneos
	}
	return nil
}

type ListaEntrenadores struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Entrenadores  []*Entrenador          `protobuf:"bytes,1,rep,name=entrenadores,proto3" json:"entrenadores,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListaEntrenadores) Reset() {
	*x = ListaEntrenadores{}
	mi := &file_comunicacion_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListaEntrenadores) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListaEntrenadores) ProtoMessage() {}

func (x *ListaEntrenadores) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListaEntrenadores.ProtoReflect.Descriptor instead.
func (*ListaEntrenadores) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{3}
}

func (x *ListaEntrenadores) GetEntrenadores() []*Entrenador {
	if x != nil {
		return x.Entrenadores
	}
	return nil
}

type Respuesta struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Mensaje       string                 `protobuf:"bytes,1,opt,name=mensaje,proto3" json:"mensaje,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Respuesta) Reset() {
	*x = Respuesta{}
	mi := &file_comunicacion_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Respuesta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Respuesta) ProtoMessage() {}

func (x *Respuesta) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Respuesta.ProtoReflect.Descriptor instead.
func (*Respuesta) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{4}
}

func (x *Respuesta) GetMensaje() string {
	if x != nil {
		return x.Mensaje
	}
	return ""
}

type InscribirRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Entrenador    *Entrenador            `protobuf:"bytes,1,opt,name=entrenador,proto3" json:"entrenador,omitempty"`
	IdTorneo      int64                  `protobuf:"varint,2,opt,name=idTorneo,proto3" json:"idTorneo,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InscribirRequest) Reset() {
	*x = InscribirRequest{}
	mi := &file_comunicacion_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InscribirRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InscribirRequest) ProtoMessage() {}

func (x *InscribirRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InscribirRequest.ProtoReflect.Descriptor instead.
func (*InscribirRequest) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{5}
}

func (x *InscribirRequest) GetEntrenador() *Entrenador {
	if x != nil {
		return x.Entrenador
	}
	return nil
}

func (x *InscribirRequest) GetIdTorneo() int64 {
	if x != nil {
		return x.IdTorneo
	}
	return 0
}

type InscribirResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Mensaje       string                 `protobuf:"bytes,1,opt,name=mensaje,proto3" json:"mensaje,omitempty"`
	En            *Entrenador            `protobuf:"bytes,2,opt,name=en,proto3" json:"en,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InscribirResponse) Reset() {
	*x = InscribirResponse{}
	mi := &file_comunicacion_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InscribirResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InscribirResponse) ProtoMessage() {}

func (x *InscribirResponse) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InscribirResponse.ProtoReflect.Descriptor instead.
func (*InscribirResponse) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{6}
}

func (x *InscribirResponse) GetMensaje() string {
	if x != nil {
		return x.Mensaje
	}
	return ""
}

func (x *InscribirResponse) GetEn() *Entrenador {
	if x != nil {
		return x.En
	}
	return nil
}

// para los combates
type Combate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IdCombate     int64                  `protobuf:"varint,1,opt,name=idCombate,proto3" json:"idCombate,omitempty"`
	IdTorneo      *Torneo                `protobuf:"bytes,2,opt,name=idTorneo,proto3" json:"idTorneo,omitempty"`
	Entrenador1   *Entrenador            `protobuf:"bytes,3,opt,name=entrenador1,proto3" json:"entrenador1,omitempty"`
	Entrenador2   *Entrenador            `protobuf:"bytes,4,opt,name=entrenador2,proto3" json:"entrenador2,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Combate) Reset() {
	*x = Combate{}
	mi := &file_comunicacion_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Combate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Combate) ProtoMessage() {}

func (x *Combate) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Combate.ProtoReflect.Descriptor instead.
func (*Combate) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{7}
}

func (x *Combate) GetIdCombate() int64 {
	if x != nil {
		return x.IdCombate
	}
	return 0
}

func (x *Combate) GetIdTorneo() *Torneo {
	if x != nil {
		return x.IdTorneo
	}
	return nil
}

func (x *Combate) GetEntrenador1() *Entrenador {
	if x != nil {
		return x.Entrenador1
	}
	return nil
}

func (x *Combate) GetEntrenador2() *Entrenador {
	if x != nil {
		return x.Entrenador2
	}
	return nil
}

type ResultadoCombate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IdCombate     int64                  `protobuf:"varint,1,opt,name=idCombate,proto3" json:"idCombate,omitempty"`
	IdGanador     int32                  `protobuf:"varint,2,opt,name=idGanador,proto3" json:"idGanador,omitempty"`
	Resumen       string                 `protobuf:"bytes,3,opt,name=resumen,proto3" json:"resumen,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResultadoCombate) Reset() {
	*x = ResultadoCombate{}
	mi := &file_comunicacion_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResultadoCombate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResultadoCombate) ProtoMessage() {}

func (x *ResultadoCombate) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResultadoCombate.ProtoReflect.Descriptor instead.
func (*ResultadoCombate) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{8}
}

func (x *ResultadoCombate) GetIdCombate() int64 {
	if x != nil {
		return x.IdCombate
	}
	return 0
}

func (x *ResultadoCombate) GetIdGanador() int32 {
	if x != nil {
		return x.IdGanador
	}
	return 0
}

func (x *ResultadoCombate) GetResumen() string {
	if x != nil {
		return x.Resumen
	}
	return ""
}

type EntrenadorID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EntrenadorID) Reset() {
	*x = EntrenadorID{}
	mi := &file_comunicacion_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EntrenadorID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntrenadorID) ProtoMessage() {}

func (x *EntrenadorID) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntrenadorID.ProtoReflect.Descriptor instead.
func (*EntrenadorID) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{9}
}

func (x *EntrenadorID) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_comunicacion_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_comunicacion_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_comunicacion_proto_rawDescGZIP(), []int{10}
}

var File_comunicacion_proto protoreflect.FileDescriptor

const file_comunicacion_proto_rawDesc = "" +
	"\n" +
	"\x12comunicacion.proto\x12\fcomunicacion\"\x9e\x01\n" +
	"\n" +
	"Entrenador\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x16\n" +
	"\x06nombre\x18\x02 \x01(\tR\x06nombre\x12\x16\n" +
	"\x06region\x18\x03 \x01(\tR\x06region\x12\x18\n" +
	"\aranking\x18\x04 \x01(\x05R\aranking\x12\x16\n" +
	"\x06estado\x18\x05 \x01(\tR\x06estado\x12\x1e\n" +
	"\n" +
	"suspension\x18\x06 \x01(\x05R\n" +
	"suspension\"\xaa\x01\n" +
	"\x06Torneo\x12\x1a\n" +
	"\bidTorneo\x18\x01 \x01(\x03R\bidTorneo\x12\x16\n" +
	"\x06region\x18\x02 \x01(\tR\x06region\x126\n" +
	"\tjugadores\x18\x03 \x03(\v2\x18.comunicacion.EntrenadorR\tjugadores\x12\x16\n" +
	"\x06estado\x18\x05 \x01(\tR\x06estado\x12\x1c\n" +
	"\tinscritos\x18\x06 \x01(\x05R\tinscritos\">\n" +
	"\fListaTorneos\x12.\n" +
	"\atorneos\x18\x01 \x03(\v2\x14.comunicacion.TorneoR\atorneos\"Q\n" +
	"\x11ListaEntrenadores\x12<\n" +
	"\fentrenadores\x18\x01 \x03(\v2\x18.comunicacion.EntrenadorR\fentrenadores\"%\n" +
	"\tRespuesta\x12\x18\n" +
	"\amensaje\x18\x01 \x01(\tR\amensaje\"h\n" +
	"\x10InscribirRequest\x128\n" +
	"\n" +
	"entrenador\x18\x01 \x01(\v2\x18.comunicacion.EntrenadorR\n" +
	"entrenador\x12\x1a\n" +
	"\bidTorneo\x18\x02 \x01(\x03R\bidTorneo\"W\n" +
	"\x11InscribirResponse\x12\x18\n" +
	"\amensaje\x18\x01 \x01(\tR\amensaje\x12(\n" +
	"\x02en\x18\x02 \x01(\v2\x18.comunicacion.EntrenadorR\x02en\"\xd1\x01\n" +
	"\aCombate\x12\x1c\n" +
	"\tidCombate\x18\x01 \x01(\x03R\tidCombate\x120\n" +
	"\bidTorneo\x18\x02 \x01(\v2\x14.comunicacion.TorneoR\bidTorneo\x12:\n" +
	"\ventrenador1\x18\x03 \x01(\v2\x18.comunicacion.EntrenadorR\ventrenador1\x12:\n" +
	"\ventrenador2\x18\x04 \x01(\v2\x18.comunicacion.EntrenadorR\ventrenador2\"h\n" +
	"\x10ResultadoCombate\x12\x1c\n" +
	"\tidCombate\x18\x01 \x01(\x03R\tidCombate\x12\x1c\n" +
	"\tidGanador\x18\x02 \x01(\x05R\tidGanador\x12\x18\n" +
	"\aresumen\x18\x03 \x01(\tR\aresumen\"\x1e\n" +
	"\fEntrenadorID\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\"\a\n" +
	"\x05empty2\x87\x03\n" +
	"\fComunicacion\x12N\n" +
	"\x12EnviarEntrenadores\x12\x1f.comunicacion.ListaEntrenadores\x1a\x17.comunicacion.Respuesta\x12@\n" +
	"\rEnviarTorneos\x12\x13.comunicacion.empty\x1a\x1a.comunicacion.ListaTorneos\x12V\n" +
	"\x13InscribirEntrenador\x12\x1e.comunicacion.InscribirRequest\x1a\x1f.comunicacion.InscribirResponse\x12@\n" +
	"\x0eAsignarCombate\x12\x15.comunicacion.Combate\x1a\x17.comunicacion.Respuesta\x12K\n" +
	"\x13ConsultarEntrenador\x12\x1a.comunicacion.EntrenadorID\x1a\x18.comunicacion.EntrenadorB\x13Z\x11grpc-server/protob\x06proto3"

var (
	file_comunicacion_proto_rawDescOnce sync.Once
	file_comunicacion_proto_rawDescData []byte
)

func file_comunicacion_proto_rawDescGZIP() []byte {
	file_comunicacion_proto_rawDescOnce.Do(func() {
		file_comunicacion_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_comunicacion_proto_rawDesc), len(file_comunicacion_proto_rawDesc)))
	})
	return file_comunicacion_proto_rawDescData
}

var file_comunicacion_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_comunicacion_proto_goTypes = []any{
	(*Entrenador)(nil),        // 0: comunicacion.Entrenador
	(*Torneo)(nil),            // 1: comunicacion.Torneo
	(*ListaTorneos)(nil),      // 2: comunicacion.ListaTorneos
	(*ListaEntrenadores)(nil), // 3: comunicacion.ListaEntrenadores
	(*Respuesta)(nil),         // 4: comunicacion.Respuesta
	(*InscribirRequest)(nil),  // 5: comunicacion.InscribirRequest
	(*InscribirResponse)(nil), // 6: comunicacion.InscribirResponse
	(*Combate)(nil),           // 7: comunicacion.Combate
	(*ResultadoCombate)(nil),  // 8: comunicacion.ResultadoCombate
	(*EntrenadorID)(nil),      // 9: comunicacion.EntrenadorID
	(*Empty)(nil),             // 10: comunicacion.empty
}
var file_comunicacion_proto_depIdxs = []int32{
	0,  // 0: comunicacion.Torneo.jugadores:type_name -> comunicacion.Entrenador
	1,  // 1: comunicacion.ListaTorneos.torneos:type_name -> comunicacion.Torneo
	0,  // 2: comunicacion.ListaEntrenadores.entrenadores:type_name -> comunicacion.Entrenador
	0,  // 3: comunicacion.InscribirRequest.entrenador:type_name -> comunicacion.Entrenador
	0,  // 4: comunicacion.InscribirResponse.en:type_name -> comunicacion.Entrenador
	1,  // 5: comunicacion.Combate.idTorneo:type_name -> comunicacion.Torneo
	0,  // 6: comunicacion.Combate.entrenador1:type_name -> comunicacion.Entrenador
	0,  // 7: comunicacion.Combate.entrenador2:type_name -> comunicacion.Entrenador
	3,  // 8: comunicacion.Comunicacion.EnviarEntrenadores:input_type -> comunicacion.ListaEntrenadores
	10, // 9: comunicacion.Comunicacion.EnviarTorneos:input_type -> comunicacion.empty
	5,  // 10: comunicacion.Comunicacion.InscribirEntrenador:input_type -> comunicacion.InscribirRequest
	7,  // 11: comunicacion.Comunicacion.AsignarCombate:input_type -> comunicacion.Combate
	9,  // 12: comunicacion.Comunicacion.ConsultarEntrenador:input_type -> comunicacion.EntrenadorID
	4,  // 13: comunicacion.Comunicacion.EnviarEntrenadores:output_type -> comunicacion.Respuesta
	2,  // 14: comunicacion.Comunicacion.EnviarTorneos:output_type -> comunicacion.ListaTorneos
	6,  // 15: comunicacion.Comunicacion.InscribirEntrenador:output_type -> comunicacion.InscribirResponse
	4,  // 16: comunicacion.Comunicacion.AsignarCombate:output_type -> comunicacion.Respuesta
	0,  // 17: comunicacion.Comunicacion.ConsultarEntrenador:output_type -> comunicacion.Entrenador
	13, // [13:18] is the sub-list for method output_type
	8,  // [8:13] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_comunicacion_proto_init() }
func file_comunicacion_proto_init() {
	if File_comunicacion_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_comunicacion_proto_rawDesc), len(file_comunicacion_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_comunicacion_proto_goTypes,
		DependencyIndexes: file_comunicacion_proto_depIdxs,
		MessageInfos:      file_comunicacion_proto_msgTypes,
	}.Build()
	File_comunicacion_proto = out.File
	file_comunicacion_proto_goTypes = nil
	file_comunicacion_proto_depIdxs = nil
}
