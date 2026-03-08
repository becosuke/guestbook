// gRPC-Gateway JSON uses camelCase field names
export interface Post {
  postId: string;
  body: string;
}

export interface ListPostsResponse {
  posts: Post[];
  nextPageToken: string;
}
