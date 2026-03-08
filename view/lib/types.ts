// gRPC-Gateway JSON uses proto field names (snake_case) with UseProtoNames: true
export interface Post {
  post_id: string;
  body: string;
}

export interface ListPostsResponse {
  posts: Post[];
  next_page_token: string;
}
