# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: search/v1/app_message.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1bsearch/v1/app_message.proto\x12\tsearch.v1\"\xa8\x01\n\x16MessageExchangeRequest\x12\x1d\n\nchannel_id\x18\x01 \x01(\tR\tchannelId\x12\x1b\n\tsender_id\x18\x02 \x01(\tR\x08senderId\x12!\n\x0crecipient_id\x18\x03 \x01(\tR\x0brecipientId\x12/\n\x07\x63ontent\x18\x04 \x01(\x0b\x32\x15.search.v1.AppMessageR\x07\x63ontent\"~\n\x0e\x41ppSendRequest\x12\x1d\n\nchannel_id\x18\x01 \x01(\tR\tchannelId\x12\x1c\n\trecipient\x18\x02 \x01(\tR\trecipient\x12/\n\x07message\x18\x03 \x01(\x0b\x32\x15.search.v1.AppMessageR\x07message\"y\n\x0f\x41ppRecvResponse\x12\x1d\n\nchannel_id\x18\x01 \x01(\tR\tchannelId\x12\x16\n\x06sender\x18\x02 \x01(\tR\x06sender\x12/\n\x07message\x18\x03 \x01(\x0b\x32\x15.search.v1.AppMessageR\x07message\"4\n\nAppMessage\x12\x12\n\x04type\x18\x01 \x01(\tR\x04type\x12\x12\n\x04\x62ody\x18\x02 \x01(\x0cR\x04\x62odyB*Z(github.com/pmontepagano/search/gen/go/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'search.v1.app_message_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z(github.com/pmontepagano/search/gen/go/v1'
  _globals['_MESSAGEEXCHANGEREQUEST']._serialized_start=43
  _globals['_MESSAGEEXCHANGEREQUEST']._serialized_end=211
  _globals['_APPSENDREQUEST']._serialized_start=213
  _globals['_APPSENDREQUEST']._serialized_end=339
  _globals['_APPRECVRESPONSE']._serialized_start=341
  _globals['_APPRECVRESPONSE']._serialized_end=462
  _globals['_APPMESSAGE']._serialized_start=464
  _globals['_APPMESSAGE']._serialized_end=516
# @@protoc_insertion_point(module_scope)