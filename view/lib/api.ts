import type { ListPostsResponse, Post } from "./pb/guestbook_pb.ts";

const API_BASE_URL = "/api/v1";

export type { ListPostsResponse, Post };

export async function getPosts(
  pageSize: number,
  pageToken: string,
): Promise<ListPostsResponse> {
  const res = await fetch(
    `${API_BASE_URL}/posts/list/${pageSize}/${pageToken}`,
  );
  if (!res.ok) throw new Error(`Failed to list posts: ${res.statusText}`);
  return res.json();
}

export async function getPost(postId: string): Promise<Post> {
  const res = await fetch(`${API_BASE_URL}/post/${postId}`);
  if (!res.ok) throw new Error(`Failed to get post: ${res.statusText}`);
  return res.json();
}

export async function createPost(body: string): Promise<Post> {
  const res = await fetch(`${API_BASE_URL}/post`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      post: { postId: "00000000-0000-0000-0000-000000000000", body },
      idempotencyKey: crypto.randomUUID(),
    }),
  });
  if (!res.ok) throw new Error(`Failed to create post: ${res.statusText}`);
  return res.json();
}

export async function updatePost(postId: string, body: string): Promise<Post> {
  const res = await fetch(`${API_BASE_URL}/post/${postId}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      post: { postId, body },
      idempotencyKey: crypto.randomUUID(),
    }),
  });
  if (!res.ok) throw new Error(`Failed to update post: ${res.statusText}`);
  return res.json();
}

export async function deletePost(postId: string): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/post/${postId}`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      postId,
      idempotencyKey: crypto.randomUUID(),
    }),
  });
  if (!res.ok) throw new Error(`Failed to delete post: ${res.statusText}`);
}
