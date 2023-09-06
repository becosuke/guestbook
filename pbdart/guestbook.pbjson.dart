//
//  Generated code. Do not modify.
//  source: guestbook.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use getPostRequestDescriptor instead')
const GetPostRequest$json = {
  '1': 'GetPostRequest',
  '2': [
    {'1': 'serial', '3': 1, '4': 1, '5': 3, '10': 'serial'},
  ],
};

/// Descriptor for `GetPostRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPostRequestDescriptor = $convert.base64Decode(
    'Cg5HZXRQb3N0UmVxdWVzdBIWCgZzZXJpYWwYASABKANSBnNlcmlhbA==');

@$core.Deprecated('Use createPostRequestDescriptor instead')
const CreatePostRequest$json = {
  '1': 'CreatePostRequest',
  '2': [
    {'1': 'post', '3': 1, '4': 1, '5': 11, '6': '.pb.Post', '8': {}, '10': 'post'},
  ],
};

/// Descriptor for `CreatePostRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPostRequestDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVQb3N0UmVxdWVzdBImCgRwb3N0GAEgASgLMggucGIuUG9zdEII+kIFigECEAFSBH'
    'Bvc3Q=');

@$core.Deprecated('Use updatePostRequestDescriptor instead')
const UpdatePostRequest$json = {
  '1': 'UpdatePostRequest',
  '2': [
    {'1': 'post', '3': 1, '4': 1, '5': 11, '6': '.pb.Post', '8': {}, '10': 'post'},
  ],
};

/// Descriptor for `UpdatePostRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePostRequestDescriptor = $convert.base64Decode(
    'ChFVcGRhdGVQb3N0UmVxdWVzdBImCgRwb3N0GAEgASgLMggucGIuUG9zdEII+kIFigECEAFSBH'
    'Bvc3Q=');

@$core.Deprecated('Use deletePostRequestDescriptor instead')
const DeletePostRequest$json = {
  '1': 'DeletePostRequest',
  '2': [
    {'1': 'serial', '3': 1, '4': 1, '5': 3, '10': 'serial'},
  ],
};

/// Descriptor for `DeletePostRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePostRequestDescriptor = $convert.base64Decode(
    'ChFEZWxldGVQb3N0UmVxdWVzdBIWCgZzZXJpYWwYASABKANSBnNlcmlhbA==');

@$core.Deprecated('Use listPostsRequestDescriptor instead')
const ListPostsRequest$json = {
  '1': 'ListPostsRequest',
  '2': [
    {'1': 'page_size', '3': 1, '4': 1, '5': 5, '8': {}, '10': 'pageSize'},
    {'1': 'page_token', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'pageToken'},
  ],
};

/// Descriptor for `ListPostsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPostsRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0UG9zdHNSZXF1ZXN0EiQKCXBhZ2Vfc2l6ZRgBIAEoBUIH+kIEGgIgAFIIcGFnZVNpem'
    'USJwoKcGFnZV90b2tlbhgCIAEoCUII+kIFcgPQAQFSCXBhZ2VUb2tlbg==');

@$core.Deprecated('Use listPostsResponseDescriptor instead')
const ListPostsResponse$json = {
  '1': 'ListPostsResponse',
  '2': [
    {'1': 'posts', '3': 1, '4': 3, '5': 11, '6': '.pb.Post', '10': 'posts'},
    {'1': 'next_page_token', '3': 2, '4': 1, '5': 9, '10': 'nextPageToken'},
  ],
};

/// Descriptor for `ListPostsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPostsResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0UG9zdHNSZXNwb25zZRIeCgVwb3N0cxgBIAMoCzIILnBiLlBvc3RSBXBvc3RzEiYKD2'
    '5leHRfcGFnZV90b2tlbhgCIAEoCVINbmV4dFBhZ2VUb2tlbg==');

@$core.Deprecated('Use postDescriptor instead')
const Post$json = {
  '1': 'Post',
  '2': [
    {'1': 'serial', '3': 1, '4': 1, '5': 3, '10': 'serial'},
    {'1': 'body', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'body'},
  ],
};

/// Descriptor for `Post`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List postDescriptor = $convert.base64Decode(
    'CgRQb3N0EhYKBnNlcmlhbBgBIAEoA1IGc2VyaWFsEh4KBGJvZHkYAiABKAlCCvpCB3IFEAEYgA'
    'FSBGJvZHk=');

