//
//  Generated code. Do not modify.
//  source: guestbook.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

class GetPostRequest extends $pb.GeneratedMessage {
  factory GetPostRequest({
    $fixnum.Int64? serial,
  }) {
    final $result = create();
    if (serial != null) {
      $result.serial = serial;
    }
    return $result;
  }
  GetPostRequest._() : super();
  factory GetPostRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetPostRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetPostRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pb'), createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'serial')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetPostRequest clone() => GetPostRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetPostRequest copyWith(void Function(GetPostRequest) updates) => super.copyWith((message) => updates(message as GetPostRequest)) as GetPostRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPostRequest create() => GetPostRequest._();
  GetPostRequest createEmptyInstance() => create();
  static $pb.PbList<GetPostRequest> createRepeated() => $pb.PbList<GetPostRequest>();
  @$core.pragma('dart2js:noInline')
  static GetPostRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetPostRequest>(create);
  static GetPostRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get serial => $_getI64(0);
  @$pb.TagNumber(1)
  set serial($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSerial() => $_has(0);
  @$pb.TagNumber(1)
  void clearSerial() => clearField(1);
}

class CreatePostRequest extends $pb.GeneratedMessage {
  factory CreatePostRequest({
    Post? post,
  }) {
    final $result = create();
    if (post != null) {
      $result.post = post;
    }
    return $result;
  }
  CreatePostRequest._() : super();
  factory CreatePostRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreatePostRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreatePostRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pb'), createEmptyInstance: create)
    ..aOM<Post>(1, _omitFieldNames ? '' : 'post', subBuilder: Post.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreatePostRequest clone() => CreatePostRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreatePostRequest copyWith(void Function(CreatePostRequest) updates) => super.copyWith((message) => updates(message as CreatePostRequest)) as CreatePostRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePostRequest create() => CreatePostRequest._();
  CreatePostRequest createEmptyInstance() => create();
  static $pb.PbList<CreatePostRequest> createRepeated() => $pb.PbList<CreatePostRequest>();
  @$core.pragma('dart2js:noInline')
  static CreatePostRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreatePostRequest>(create);
  static CreatePostRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Post get post => $_getN(0);
  @$pb.TagNumber(1)
  set post(Post v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasPost() => $_has(0);
  @$pb.TagNumber(1)
  void clearPost() => clearField(1);
  @$pb.TagNumber(1)
  Post ensurePost() => $_ensure(0);
}

class UpdatePostRequest extends $pb.GeneratedMessage {
  factory UpdatePostRequest({
    Post? post,
  }) {
    final $result = create();
    if (post != null) {
      $result.post = post;
    }
    return $result;
  }
  UpdatePostRequest._() : super();
  factory UpdatePostRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UpdatePostRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdatePostRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pb'), createEmptyInstance: create)
    ..aOM<Post>(1, _omitFieldNames ? '' : 'post', subBuilder: Post.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UpdatePostRequest clone() => UpdatePostRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UpdatePostRequest copyWith(void Function(UpdatePostRequest) updates) => super.copyWith((message) => updates(message as UpdatePostRequest)) as UpdatePostRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdatePostRequest create() => UpdatePostRequest._();
  UpdatePostRequest createEmptyInstance() => create();
  static $pb.PbList<UpdatePostRequest> createRepeated() => $pb.PbList<UpdatePostRequest>();
  @$core.pragma('dart2js:noInline')
  static UpdatePostRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdatePostRequest>(create);
  static UpdatePostRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Post get post => $_getN(0);
  @$pb.TagNumber(1)
  set post(Post v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasPost() => $_has(0);
  @$pb.TagNumber(1)
  void clearPost() => clearField(1);
  @$pb.TagNumber(1)
  Post ensurePost() => $_ensure(0);
}

class DeletePostRequest extends $pb.GeneratedMessage {
  factory DeletePostRequest({
    $fixnum.Int64? serial,
  }) {
    final $result = create();
    if (serial != null) {
      $result.serial = serial;
    }
    return $result;
  }
  DeletePostRequest._() : super();
  factory DeletePostRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DeletePostRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeletePostRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pb'), createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'serial')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DeletePostRequest clone() => DeletePostRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DeletePostRequest copyWith(void Function(DeletePostRequest) updates) => super.copyWith((message) => updates(message as DeletePostRequest)) as DeletePostRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeletePostRequest create() => DeletePostRequest._();
  DeletePostRequest createEmptyInstance() => create();
  static $pb.PbList<DeletePostRequest> createRepeated() => $pb.PbList<DeletePostRequest>();
  @$core.pragma('dart2js:noInline')
  static DeletePostRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeletePostRequest>(create);
  static DeletePostRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get serial => $_getI64(0);
  @$pb.TagNumber(1)
  set serial($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSerial() => $_has(0);
  @$pb.TagNumber(1)
  void clearSerial() => clearField(1);
}

class ListPostsRequest extends $pb.GeneratedMessage {
  factory ListPostsRequest({
    $core.int? pageSize,
    $core.String? pageToken,
  }) {
    final $result = create();
    if (pageSize != null) {
      $result.pageSize = pageSize;
    }
    if (pageToken != null) {
      $result.pageToken = pageToken;
    }
    return $result;
  }
  ListPostsRequest._() : super();
  factory ListPostsRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListPostsRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListPostsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pb'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'pageSize', $pb.PbFieldType.O3)
    ..aOS(2, _omitFieldNames ? '' : 'pageToken')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListPostsRequest clone() => ListPostsRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListPostsRequest copyWith(void Function(ListPostsRequest) updates) => super.copyWith((message) => updates(message as ListPostsRequest)) as ListPostsRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPostsRequest create() => ListPostsRequest._();
  ListPostsRequest createEmptyInstance() => create();
  static $pb.PbList<ListPostsRequest> createRepeated() => $pb.PbList<ListPostsRequest>();
  @$core.pragma('dart2js:noInline')
  static ListPostsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListPostsRequest>(create);
  static ListPostsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get pageSize => $_getIZ(0);
  @$pb.TagNumber(1)
  set pageSize($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPageSize() => $_has(0);
  @$pb.TagNumber(1)
  void clearPageSize() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get pageToken => $_getSZ(1);
  @$pb.TagNumber(2)
  set pageToken($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPageToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearPageToken() => clearField(2);
}

class ListPostsResponse extends $pb.GeneratedMessage {
  factory ListPostsResponse({
    $core.Iterable<Post>? posts,
    $core.String? nextPageToken,
  }) {
    final $result = create();
    if (posts != null) {
      $result.posts.addAll(posts);
    }
    if (nextPageToken != null) {
      $result.nextPageToken = nextPageToken;
    }
    return $result;
  }
  ListPostsResponse._() : super();
  factory ListPostsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListPostsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListPostsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pb'), createEmptyInstance: create)
    ..pc<Post>(1, _omitFieldNames ? '' : 'posts', $pb.PbFieldType.PM, subBuilder: Post.create)
    ..aOS(2, _omitFieldNames ? '' : 'nextPageToken')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListPostsResponse clone() => ListPostsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListPostsResponse copyWith(void Function(ListPostsResponse) updates) => super.copyWith((message) => updates(message as ListPostsResponse)) as ListPostsResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPostsResponse create() => ListPostsResponse._();
  ListPostsResponse createEmptyInstance() => create();
  static $pb.PbList<ListPostsResponse> createRepeated() => $pb.PbList<ListPostsResponse>();
  @$core.pragma('dart2js:noInline')
  static ListPostsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListPostsResponse>(create);
  static ListPostsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<Post> get posts => $_getList(0);

  @$pb.TagNumber(2)
  $core.String get nextPageToken => $_getSZ(1);
  @$pb.TagNumber(2)
  set nextPageToken($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasNextPageToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearNextPageToken() => clearField(2);
}

class Post extends $pb.GeneratedMessage {
  factory Post({
    $fixnum.Int64? serial,
    $core.String? body,
  }) {
    final $result = create();
    if (serial != null) {
      $result.serial = serial;
    }
    if (body != null) {
      $result.body = body;
    }
    return $result;
  }
  Post._() : super();
  factory Post.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Post.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Post', package: const $pb.PackageName(_omitMessageNames ? '' : 'pb'), createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'serial')
    ..aOS(2, _omitFieldNames ? '' : 'body')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Post clone() => Post()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Post copyWith(void Function(Post) updates) => super.copyWith((message) => updates(message as Post)) as Post;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Post create() => Post._();
  Post createEmptyInstance() => create();
  static $pb.PbList<Post> createRepeated() => $pb.PbList<Post>();
  @$core.pragma('dart2js:noInline')
  static Post getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Post>(create);
  static Post? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get serial => $_getI64(0);
  @$pb.TagNumber(1)
  set serial($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSerial() => $_has(0);
  @$pb.TagNumber(1)
  void clearSerial() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get body => $_getSZ(1);
  @$pb.TagNumber(2)
  set body($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasBody() => $_has(1);
  @$pb.TagNumber(2)
  void clearBody() => clearField(2);
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
