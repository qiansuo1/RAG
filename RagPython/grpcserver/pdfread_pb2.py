# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: pdfread.proto
# Protobuf Python Version: 5.27.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    27,
    2,
    '',
    'pdfread.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\rpdfread.proto\x12\npdfservice\"\x1f\n\nPdfRequest\x12\x11\n\tfile_path\x18\x01 \x01(\t\"\x9c\x01\n\x0bPdfResponse\x12\x13\n\x0bpage_number\x18\x01 \x01(\x05\x12\x16\n\x0esentence_chunk\x18\x02 \x01(\t\x12\x18\n\x10\x63hunk_char_count\x18\x03 \x01(\x05\x12\x18\n\x10\x63hunk_word_count\x18\x04 \x01(\x05\x12\x19\n\x11\x63hunk_token_count\x18\x05 \x01(\x05\x12\x11\n\tembedding\x18\x06 \x03(\x02\" \n\x10VectorizeRequest\x12\x0c\n\x04text\x18\x01 \x01(\t\"#\n\x11VectorizeResponse\x12\x0e\n\x06vector\x18\x01 \x03(\x02\x32\xa0\x01\n\nPdfService\x12\x42\n\x0b\x45xtractText\x12\x16.pdfservice.PdfRequest\x1a\x17.pdfservice.PdfResponse\"\x00\x30\x01\x12N\n\rVectorizeText\x12\x1c.pdfservice.VectorizeRequest\x1a\x1d.pdfservice.VectorizeResponse\"\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'pdfread_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_PDFREQUEST']._serialized_start=29
  _globals['_PDFREQUEST']._serialized_end=60
  _globals['_PDFRESPONSE']._serialized_start=63
  _globals['_PDFRESPONSE']._serialized_end=219
  _globals['_VECTORIZEREQUEST']._serialized_start=221
  _globals['_VECTORIZEREQUEST']._serialized_end=253
  _globals['_VECTORIZERESPONSE']._serialized_start=255
  _globals['_VECTORIZERESPONSE']._serialized_end=290
  _globals['_PDFSERVICE']._serialized_start=293
  _globals['_PDFSERVICE']._serialized_end=453
# @@protoc_insertion_point(module_scope)
