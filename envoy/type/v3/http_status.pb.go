// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: envoy/type/v3/http_status.proto

package typev3

import (
	_ "github.com/cncf/xds/go/udpa/annotations"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// HTTP response codes supported in Envoy.
// For more details: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
type StatusCode int32

const (
	// Empty - This code not part of the HTTP status code specification, but it is needed for proto
	// `enum` type.
	StatusCode_Empty                         StatusCode = 0
	StatusCode_Continue                      StatusCode = 100
	StatusCode_OK                            StatusCode = 200
	StatusCode_Created                       StatusCode = 201
	StatusCode_Accepted                      StatusCode = 202
	StatusCode_NonAuthoritativeInformation   StatusCode = 203
	StatusCode_NoContent                     StatusCode = 204
	StatusCode_ResetContent                  StatusCode = 205
	StatusCode_PartialContent                StatusCode = 206
	StatusCode_MultiStatus                   StatusCode = 207
	StatusCode_AlreadyReported               StatusCode = 208
	StatusCode_IMUsed                        StatusCode = 226
	StatusCode_MultipleChoices               StatusCode = 300
	StatusCode_MovedPermanently              StatusCode = 301
	StatusCode_Found                         StatusCode = 302
	StatusCode_SeeOther                      StatusCode = 303
	StatusCode_NotModified                   StatusCode = 304
	StatusCode_UseProxy                      StatusCode = 305
	StatusCode_TemporaryRedirect             StatusCode = 307
	StatusCode_PermanentRedirect             StatusCode = 308
	StatusCode_BadRequest                    StatusCode = 400
	StatusCode_Unauthorized                  StatusCode = 401
	StatusCode_PaymentRequired               StatusCode = 402
	StatusCode_Forbidden                     StatusCode = 403
	StatusCode_NotFound                      StatusCode = 404
	StatusCode_MethodNotAllowed              StatusCode = 405
	StatusCode_NotAcceptable                 StatusCode = 406
	StatusCode_ProxyAuthenticationRequired   StatusCode = 407
	StatusCode_RequestTimeout                StatusCode = 408
	StatusCode_Conflict                      StatusCode = 409
	StatusCode_Gone                          StatusCode = 410
	StatusCode_LengthRequired                StatusCode = 411
	StatusCode_PreconditionFailed            StatusCode = 412
	StatusCode_PayloadTooLarge               StatusCode = 413
	StatusCode_URITooLong                    StatusCode = 414
	StatusCode_UnsupportedMediaType          StatusCode = 415
	StatusCode_RangeNotSatisfiable           StatusCode = 416
	StatusCode_ExpectationFailed             StatusCode = 417
	StatusCode_MisdirectedRequest            StatusCode = 421
	StatusCode_UnprocessableEntity           StatusCode = 422
	StatusCode_Locked                        StatusCode = 423
	StatusCode_FailedDependency              StatusCode = 424
	StatusCode_UpgradeRequired               StatusCode = 426
	StatusCode_PreconditionRequired          StatusCode = 428
	StatusCode_TooManyRequests               StatusCode = 429
	StatusCode_RequestHeaderFieldsTooLarge   StatusCode = 431
	StatusCode_InternalServerError           StatusCode = 500
	StatusCode_NotImplemented                StatusCode = 501
	StatusCode_BadGateway                    StatusCode = 502
	StatusCode_ServiceUnavailable            StatusCode = 503
	StatusCode_GatewayTimeout                StatusCode = 504
	StatusCode_HTTPVersionNotSupported       StatusCode = 505
	StatusCode_VariantAlsoNegotiates         StatusCode = 506
	StatusCode_InsufficientStorage           StatusCode = 507
	StatusCode_LoopDetected                  StatusCode = 508
	StatusCode_NotExtended                   StatusCode = 510
	StatusCode_NetworkAuthenticationRequired StatusCode = 511
)

// Enum value maps for StatusCode.
var (
	StatusCode_name = map[int32]string{
		0:   "Empty",
		100: "Continue",
		200: "OK",
		201: "Created",
		202: "Accepted",
		203: "NonAuthoritativeInformation",
		204: "NoContent",
		205: "ResetContent",
		206: "PartialContent",
		207: "MultiStatus",
		208: "AlreadyReported",
		226: "IMUsed",
		300: "MultipleChoices",
		301: "MovedPermanently",
		302: "Found",
		303: "SeeOther",
		304: "NotModified",
		305: "UseProxy",
		307: "TemporaryRedirect",
		308: "PermanentRedirect",
		400: "BadRequest",
		401: "Unauthorized",
		402: "PaymentRequired",
		403: "Forbidden",
		404: "NotFound",
		405: "MethodNotAllowed",
		406: "NotAcceptable",
		407: "ProxyAuthenticationRequired",
		408: "RequestTimeout",
		409: "Conflict",
		410: "Gone",
		411: "LengthRequired",
		412: "PreconditionFailed",
		413: "PayloadTooLarge",
		414: "URITooLong",
		415: "UnsupportedMediaType",
		416: "RangeNotSatisfiable",
		417: "ExpectationFailed",
		421: "MisdirectedRequest",
		422: "UnprocessableEntity",
		423: "Locked",
		424: "FailedDependency",
		426: "UpgradeRequired",
		428: "PreconditionRequired",
		429: "TooManyRequests",
		431: "RequestHeaderFieldsTooLarge",
		500: "InternalServerError",
		501: "NotImplemented",
		502: "BadGateway",
		503: "ServiceUnavailable",
		504: "GatewayTimeout",
		505: "HTTPVersionNotSupported",
		506: "VariantAlsoNegotiates",
		507: "InsufficientStorage",
		508: "LoopDetected",
		510: "NotExtended",
		511: "NetworkAuthenticationRequired",
	}
	StatusCode_value = map[string]int32{
		"Empty":                         0,
		"Continue":                      100,
		"OK":                            200,
		"Created":                       201,
		"Accepted":                      202,
		"NonAuthoritativeInformation":   203,
		"NoContent":                     204,
		"ResetContent":                  205,
		"PartialContent":                206,
		"MultiStatus":                   207,
		"AlreadyReported":               208,
		"IMUsed":                        226,
		"MultipleChoices":               300,
		"MovedPermanently":              301,
		"Found":                         302,
		"SeeOther":                      303,
		"NotModified":                   304,
		"UseProxy":                      305,
		"TemporaryRedirect":             307,
		"PermanentRedirect":             308,
		"BadRequest":                    400,
		"Unauthorized":                  401,
		"PaymentRequired":               402,
		"Forbidden":                     403,
		"NotFound":                      404,
		"MethodNotAllowed":              405,
		"NotAcceptable":                 406,
		"ProxyAuthenticationRequired":   407,
		"RequestTimeout":                408,
		"Conflict":                      409,
		"Gone":                          410,
		"LengthRequired":                411,
		"PreconditionFailed":            412,
		"PayloadTooLarge":               413,
		"URITooLong":                    414,
		"UnsupportedMediaType":          415,
		"RangeNotSatisfiable":           416,
		"ExpectationFailed":             417,
		"MisdirectedRequest":            421,
		"UnprocessableEntity":           422,
		"Locked":                        423,
		"FailedDependency":              424,
		"UpgradeRequired":               426,
		"PreconditionRequired":          428,
		"TooManyRequests":               429,
		"RequestHeaderFieldsTooLarge":   431,
		"InternalServerError":           500,
		"NotImplemented":                501,
		"BadGateway":                    502,
		"ServiceUnavailable":            503,
		"GatewayTimeout":                504,
		"HTTPVersionNotSupported":       505,
		"VariantAlsoNegotiates":         506,
		"InsufficientStorage":           507,
		"LoopDetected":                  508,
		"NotExtended":                   510,
		"NetworkAuthenticationRequired": 511,
	}
)

func (x StatusCode) Enum() *StatusCode {
	p := new(StatusCode)
	*p = x
	return p
}

func (x StatusCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusCode) Descriptor() protoreflect.EnumDescriptor {
	return file_envoy_type_v3_http_status_proto_enumTypes[0].Descriptor()
}

func (StatusCode) Type() protoreflect.EnumType {
	return &file_envoy_type_v3_http_status_proto_enumTypes[0]
}

func (x StatusCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusCode.Descriptor instead.
func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return file_envoy_type_v3_http_status_proto_rawDescGZIP(), []int{0}
}

// HTTP status.
type HttpStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Supplies HTTP response code.
	Code StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=envoy.type.v3.StatusCode" json:"code,omitempty"`
}

func (x *HttpStatus) Reset() {
	*x = HttpStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_type_v3_http_status_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpStatus) ProtoMessage() {}

func (x *HttpStatus) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_type_v3_http_status_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpStatus.ProtoReflect.Descriptor instead.
func (*HttpStatus) Descriptor() ([]byte, []int) {
	return file_envoy_type_v3_http_status_proto_rawDescGZIP(), []int{0}
}

func (x *HttpStatus) GetCode() StatusCode {
	if x != nil {
		return x.Code
	}
	return StatusCode_Empty
}

var File_envoy_type_v3_http_status_proto protoreflect.FileDescriptor

var file_envoy_type_v3_http_status_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x33, 0x2f,
	0x68, 0x74, 0x74, 0x70, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0d, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x33,
	0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x21, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a, 0x0a, 0x48,
	0x74, 0x74, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x33, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f,
	0x64, 0x65, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x82, 0x01, 0x04, 0x10, 0x01, 0x20, 0x00, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x3a, 0x1c, 0x9a, 0xc5, 0x88, 0x1e, 0x17, 0x0a, 0x15, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x2a, 0xb5, 0x09, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08,
	0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x10, 0x64, 0x12, 0x07, 0x0a, 0x02, 0x4f, 0x4b,
	0x10, 0xc8, 0x01, 0x12, 0x0c, 0x0a, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x10, 0xc9,
	0x01, 0x12, 0x0d, 0x0a, 0x08, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x10, 0xca, 0x01,
	0x12, 0x20, 0x0a, 0x1b, 0x4e, 0x6f, 0x6e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10,
	0xcb, 0x01, 0x12, 0x0e, 0x0a, 0x09, 0x4e, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x10,
	0xcc, 0x01, 0x12, 0x11, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x10, 0xcd, 0x01, 0x12, 0x13, 0x0a, 0x0e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x10, 0xce, 0x01, 0x12, 0x10, 0x0a, 0x0b, 0x4d, 0x75,
	0x6c, 0x74, 0x69, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x10, 0xcf, 0x01, 0x12, 0x14, 0x0a, 0x0f,
	0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x10,
	0xd0, 0x01, 0x12, 0x0b, 0x0a, 0x06, 0x49, 0x4d, 0x55, 0x73, 0x65, 0x64, 0x10, 0xe2, 0x01, 0x12,
	0x14, 0x0a, 0x0f, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x65, 0x43, 0x68, 0x6f, 0x69, 0x63,
	0x65, 0x73, 0x10, 0xac, 0x02, 0x12, 0x15, 0x0a, 0x10, 0x4d, 0x6f, 0x76, 0x65, 0x64, 0x50, 0x65,
	0x72, 0x6d, 0x61, 0x6e, 0x65, 0x6e, 0x74, 0x6c, 0x79, 0x10, 0xad, 0x02, 0x12, 0x0a, 0x0a, 0x05,
	0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xae, 0x02, 0x12, 0x0d, 0x0a, 0x08, 0x53, 0x65, 0x65, 0x4f,
	0x74, 0x68, 0x65, 0x72, 0x10, 0xaf, 0x02, 0x12, 0x10, 0x0a, 0x0b, 0x4e, 0x6f, 0x74, 0x4d, 0x6f,
	0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x10, 0xb0, 0x02, 0x12, 0x0d, 0x0a, 0x08, 0x55, 0x73, 0x65,
	0x50, 0x72, 0x6f, 0x78, 0x79, 0x10, 0xb1, 0x02, 0x12, 0x16, 0x0a, 0x11, 0x54, 0x65, 0x6d, 0x70,
	0x6f, 0x72, 0x61, 0x72, 0x79, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x10, 0xb3, 0x02,
	0x12, 0x16, 0x0a, 0x11, 0x50, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x10, 0xb4, 0x02, 0x12, 0x0f, 0x0a, 0x0a, 0x42, 0x61, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x90, 0x03, 0x12, 0x11, 0x0a, 0x0c, 0x55, 0x6e, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x10, 0x91, 0x03, 0x12, 0x14, 0x0a, 0x0f,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x10,
	0x92, 0x03, 0x12, 0x0e, 0x0a, 0x09, 0x46, 0x6f, 0x72, 0x62, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x10,
	0x93, 0x03, 0x12, 0x0d, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0x94,
	0x03, 0x12, 0x15, 0x0a, 0x10, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x6f, 0x74, 0x41, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x64, 0x10, 0x95, 0x03, 0x12, 0x12, 0x0a, 0x0d, 0x4e, 0x6f, 0x74, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x10, 0x96, 0x03, 0x12, 0x20, 0x0a, 0x1b,
	0x50, 0x72, 0x6f, 0x78, 0x79, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x10, 0x97, 0x03, 0x12, 0x13,
	0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x10, 0x98, 0x03, 0x12, 0x0d, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x10,
	0x99, 0x03, 0x12, 0x09, 0x0a, 0x04, 0x47, 0x6f, 0x6e, 0x65, 0x10, 0x9a, 0x03, 0x12, 0x13, 0x0a,
	0x0e, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x10,
	0x9b, 0x03, 0x12, 0x17, 0x0a, 0x12, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x9c, 0x03, 0x12, 0x14, 0x0a, 0x0f, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x6f, 0x6f, 0x4c, 0x61, 0x72, 0x67, 0x65, 0x10, 0x9d,
	0x03, 0x12, 0x0f, 0x0a, 0x0a, 0x55, 0x52, 0x49, 0x54, 0x6f, 0x6f, 0x4c, 0x6f, 0x6e, 0x67, 0x10,
	0x9e, 0x03, 0x12, 0x19, 0x0a, 0x14, 0x55, 0x6e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x10, 0x9f, 0x03, 0x12, 0x18, 0x0a,
	0x13, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x4e, 0x6f, 0x74, 0x53, 0x61, 0x74, 0x69, 0x73, 0x66, 0x69,
	0x61, 0x62, 0x6c, 0x65, 0x10, 0xa0, 0x03, 0x12, 0x16, 0x0a, 0x11, 0x45, 0x78, 0x70, 0x65, 0x63,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0xa1, 0x03, 0x12,
	0x17, 0x0a, 0x12, 0x4d, 0x69, 0x73, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0xa5, 0x03, 0x12, 0x18, 0x0a, 0x13, 0x55, 0x6e, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x10,
	0xa6, 0x03, 0x12, 0x0b, 0x0a, 0x06, 0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x10, 0xa7, 0x03, 0x12,
	0x15, 0x0a, 0x10, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65,
	0x6e, 0x63, 0x79, 0x10, 0xa8, 0x03, 0x12, 0x14, 0x0a, 0x0f, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x10, 0xaa, 0x03, 0x12, 0x19, 0x0a, 0x14,
	0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x69, 0x72, 0x65, 0x64, 0x10, 0xac, 0x03, 0x12, 0x14, 0x0a, 0x0f, 0x54, 0x6f, 0x6f, 0x4d, 0x61,
	0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x10, 0xad, 0x03, 0x12, 0x20, 0x0a,
	0x1b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x54, 0x6f, 0x6f, 0x4c, 0x61, 0x72, 0x67, 0x65, 0x10, 0xaf, 0x03, 0x12,
	0x18, 0x0a, 0x13, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0xf4, 0x03, 0x12, 0x13, 0x0a, 0x0e, 0x4e, 0x6f, 0x74,
	0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x65, 0x64, 0x10, 0xf5, 0x03, 0x12, 0x0f,
	0x0a, 0x0a, 0x42, 0x61, 0x64, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x10, 0xf6, 0x03, 0x12,
	0x17, 0x0a, 0x12, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x55, 0x6e, 0x61, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x10, 0xf7, 0x03, 0x12, 0x13, 0x0a, 0x0e, 0x47, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x10, 0xf8, 0x03, 0x12, 0x1c, 0x0a,
	0x17, 0x48, 0x54, 0x54, 0x50, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x74, 0x53,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x10, 0xf9, 0x03, 0x12, 0x1a, 0x0a, 0x15, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x41, 0x6c, 0x73, 0x6f, 0x4e, 0x65, 0x67, 0x6f, 0x74, 0x69,
	0x61, 0x74, 0x65, 0x73, 0x10, 0xfa, 0x03, 0x12, 0x18, 0x0a, 0x13, 0x49, 0x6e, 0x73, 0x75, 0x66,
	0x66, 0x69, 0x63, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x10, 0xfb,
	0x03, 0x12, 0x11, 0x0a, 0x0c, 0x4c, 0x6f, 0x6f, 0x70, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x65,
	0x64, 0x10, 0xfc, 0x03, 0x12, 0x10, 0x0a, 0x0b, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x74, 0x65, 0x6e,
	0x64, 0x65, 0x64, 0x10, 0xfe, 0x03, 0x12, 0x22, 0x0a, 0x1d, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x10, 0xff, 0x03, 0x42, 0x75, 0xba, 0x80, 0xc8, 0xd1,
	0x06, 0x02, 0x10, 0x02, 0x0a, 0x1b, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72,
	0x6f, 0x78, 0x79, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76,
	0x33, 0x42, 0x0f, 0x48, 0x74, 0x74, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x33, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x76,
	0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_envoy_type_v3_http_status_proto_rawDescOnce sync.Once
	file_envoy_type_v3_http_status_proto_rawDescData = file_envoy_type_v3_http_status_proto_rawDesc
)

func file_envoy_type_v3_http_status_proto_rawDescGZIP() []byte {
	file_envoy_type_v3_http_status_proto_rawDescOnce.Do(func() {
		file_envoy_type_v3_http_status_proto_rawDescData = protoimpl.X.CompressGZIP(file_envoy_type_v3_http_status_proto_rawDescData)
	})
	return file_envoy_type_v3_http_status_proto_rawDescData
}

var file_envoy_type_v3_http_status_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_envoy_type_v3_http_status_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_envoy_type_v3_http_status_proto_goTypes = []interface{}{
	(StatusCode)(0),    // 0: envoy.type.v3.StatusCode
	(*HttpStatus)(nil), // 1: envoy.type.v3.HttpStatus
}
var file_envoy_type_v3_http_status_proto_depIdxs = []int32{
	0, // 0: envoy.type.v3.HttpStatus.code:type_name -> envoy.type.v3.StatusCode
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_envoy_type_v3_http_status_proto_init() }
func file_envoy_type_v3_http_status_proto_init() {
	if File_envoy_type_v3_http_status_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envoy_type_v3_http_status_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_envoy_type_v3_http_status_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envoy_type_v3_http_status_proto_goTypes,
		DependencyIndexes: file_envoy_type_v3_http_status_proto_depIdxs,
		EnumInfos:         file_envoy_type_v3_http_status_proto_enumTypes,
		MessageInfos:      file_envoy_type_v3_http_status_proto_msgTypes,
	}.Build()
	File_envoy_type_v3_http_status_proto = out.File
	file_envoy_type_v3_http_status_proto_rawDesc = nil
	file_envoy_type_v3_http_status_proto_goTypes = nil
	file_envoy_type_v3_http_status_proto_depIdxs = nil
}
