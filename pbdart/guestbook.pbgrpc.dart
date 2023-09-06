//
//  Generated code. Do not modify.
//  source: guestbook.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'google/protobuf/empty.pb.dart' as $1;
import 'guestbook.pb.dart' as $0;

export 'guestbook.pb.dart';

@$pb.GrpcServiceName('pb.GuestbookService')
class GuestbookServiceClient extends $grpc.Client {
  static final _$getPost = $grpc.ClientMethod<$0.GetPostRequest, $0.Post>(
      '/pb.GuestbookService/GetPost',
      ($0.GetPostRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.Post.fromBuffer(value));
  static final _$createPost = $grpc.ClientMethod<$0.CreatePostRequest, $0.Post>(
      '/pb.GuestbookService/CreatePost',
      ($0.CreatePostRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.Post.fromBuffer(value));
  static final _$updatePost = $grpc.ClientMethod<$0.UpdatePostRequest, $0.Post>(
      '/pb.GuestbookService/UpdatePost',
      ($0.UpdatePostRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.Post.fromBuffer(value));
  static final _$deletePost = $grpc.ClientMethod<$0.DeletePostRequest, $1.Empty>(
      '/pb.GuestbookService/DeletePost',
      ($0.DeletePostRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $1.Empty.fromBuffer(value));
  static final _$listPosts = $grpc.ClientMethod<$0.ListPostsRequest, $0.ListPostsResponse>(
      '/pb.GuestbookService/ListPosts',
      ($0.ListPostsRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.ListPostsResponse.fromBuffer(value));

  GuestbookServiceClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options,
        interceptors: interceptors);

  $grpc.ResponseFuture<$0.Post> getPost($0.GetPostRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$getPost, request, options: options);
  }

  $grpc.ResponseFuture<$0.Post> createPost($0.CreatePostRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$createPost, request, options: options);
  }

  $grpc.ResponseFuture<$0.Post> updatePost($0.UpdatePostRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$updatePost, request, options: options);
  }

  $grpc.ResponseFuture<$1.Empty> deletePost($0.DeletePostRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$deletePost, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListPostsResponse> listPosts($0.ListPostsRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$listPosts, request, options: options);
  }
}

@$pb.GrpcServiceName('pb.GuestbookService')
abstract class GuestbookServiceBase extends $grpc.Service {
  $core.String get $name => 'pb.GuestbookService';

  GuestbookServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetPostRequest, $0.Post>(
        'GetPost',
        getPost_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetPostRequest.fromBuffer(value),
        ($0.Post value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePostRequest, $0.Post>(
        'CreatePost',
        createPost_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreatePostRequest.fromBuffer(value),
        ($0.Post value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdatePostRequest, $0.Post>(
        'UpdatePost',
        updatePost_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdatePostRequest.fromBuffer(value),
        ($0.Post value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePostRequest, $1.Empty>(
        'DeletePost',
        deletePost_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeletePostRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPostsRequest, $0.ListPostsResponse>(
        'ListPosts',
        listPosts_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListPostsRequest.fromBuffer(value),
        ($0.ListPostsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.Post> getPost_Pre($grpc.ServiceCall call, $async.Future<$0.GetPostRequest> request) async {
    return getPost(call, await request);
  }

  $async.Future<$0.Post> createPost_Pre($grpc.ServiceCall call, $async.Future<$0.CreatePostRequest> request) async {
    return createPost(call, await request);
  }

  $async.Future<$0.Post> updatePost_Pre($grpc.ServiceCall call, $async.Future<$0.UpdatePostRequest> request) async {
    return updatePost(call, await request);
  }

  $async.Future<$1.Empty> deletePost_Pre($grpc.ServiceCall call, $async.Future<$0.DeletePostRequest> request) async {
    return deletePost(call, await request);
  }

  $async.Future<$0.ListPostsResponse> listPosts_Pre($grpc.ServiceCall call, $async.Future<$0.ListPostsRequest> request) async {
    return listPosts(call, await request);
  }

  $async.Future<$0.Post> getPost($grpc.ServiceCall call, $0.GetPostRequest request);
  $async.Future<$0.Post> createPost($grpc.ServiceCall call, $0.CreatePostRequest request);
  $async.Future<$0.Post> updatePost($grpc.ServiceCall call, $0.UpdatePostRequest request);
  $async.Future<$1.Empty> deletePost($grpc.ServiceCall call, $0.DeletePostRequest request);
  $async.Future<$0.ListPostsResponse> listPosts($grpc.ServiceCall call, $0.ListPostsRequest request);
}
