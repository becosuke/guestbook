import { useSignal } from "@preact/signals";
import type { Post } from "../lib/types.ts";
import { createPost, deletePost, getPosts, updatePost } from "../lib/api.ts";

const PAGE_SIZE = 10;

export default function PostManager() {
  const posts = useSignal<Post[]>([]);
  const nextPageToken = useSignal("");
  const newBody = useSignal("");
  const editingId = useSignal<string | null>(null);
  const editingBody = useSignal("");
  const loading = useSignal(false);
  const error = useSignal("");
  const initialized = useSignal(false);

  async function loadPosts(pageToken = "") {
    loading.value = true;
    error.value = "";
    try {
      const res = await getPosts(PAGE_SIZE, pageToken);
      posts.value = pageToken
        ? [...posts.value, ...res.posts]
        : res.posts;
      nextPageToken.value = res.nextPageToken;
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to load posts";
    } finally {
      loading.value = false;
    }
  }

  if (!initialized.value) {
    initialized.value = true;
    loadPosts();
  }

  async function handleCreate() {
    const body = newBody.value.trim();
    if (!body) return;
    loading.value = true;
    error.value = "";
    try {
      await createPost(body);
      newBody.value = "";
      posts.value = [];
      nextPageToken.value = "";
      await loadPosts();
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to create post";
    } finally {
      loading.value = false;
    }
  }

  async function handleUpdate(postId: string) {
    const body = editingBody.value.trim();
    if (!body) return;
    loading.value = true;
    error.value = "";
    try {
      const updated = await updatePost(postId, body);
      posts.value = posts.value.map((p) =>
        p.postId === postId ? updated : p
      );
      editingId.value = null;
      editingBody.value = "";
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to update post";
    } finally {
      loading.value = false;
    }
  }

  async function handleDelete(postId: string) {
    loading.value = true;
    error.value = "";
    try {
      await deletePost(postId);
      posts.value = posts.value.filter((p) => p.postId !== postId);
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to delete post";
    } finally {
      loading.value = false;
    }
  }

  function startEdit(post: Post) {
    editingId.value = post.postId;
    editingBody.value = post.body;
  }

  function cancelEdit() {
    editingId.value = null;
    editingBody.value = "";
  }

  return (
    <div class="post-manager">
      <div class="create-form">
        <input
          type="text"
          class="text-input"
          placeholder="Write a new post..."
          value={newBody.value}
          onInput={(e) =>
            newBody.value = (e.target as HTMLInputElement).value}
          onKeyDown={(e) => {
            if (e.key === "Enter") handleCreate();
          }}
          maxLength={128}
        />
        <button
          class="btn btn-primary"
          onClick={handleCreate}
          disabled={loading.value || !newBody.value.trim()}
        >
          Post
        </button>
      </div>

      {error.value && <div class="error-message">{error.value}</div>}

      <div class="post-list">
        {posts.value.map((post) => (
          <div class="post-card" key={post.postId}>
            {editingId.value === post.postId
              ? (
                <div class="edit-form">
                  <input
                    type="text"
                    class="text-input"
                    value={editingBody.value}
                    onInput={(e) =>
                      editingBody.value =
                        (e.target as HTMLInputElement).value}
                    onKeyDown={(e) => {
                      if (e.key === "Enter") handleUpdate(post.postId);
                      if (e.key === "Escape") cancelEdit();
                    }}
                    maxLength={128}
                  />
                  <button
                    class="btn btn-primary"
                    onClick={() => handleUpdate(post.postId)}
                    disabled={loading.value || !editingBody.value.trim()}
                  >
                    Save
                  </button>
                  <button class="btn btn-secondary" onClick={cancelEdit}>
                    Cancel
                  </button>
                </div>
              )
              : (
                <div class="post-content">
                  <p class="post-body">{post.body}</p>
                  <div class="post-actions">
                    <button
                      class="btn btn-small"
                      onClick={() => startEdit(post)}
                    >
                      Edit
                    </button>
                    <button
                      class="btn btn-small btn-danger"
                      onClick={() => handleDelete(post.postId)}
                      disabled={loading.value}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              )}
            <div class="post-id">{post.postId}</div>
          </div>
        ))}
      </div>

      {posts.value.length === 0 && !loading.value && (
        <p class="empty-message">No posts yet. Create one above!</p>
      )}

      {loading.value && <p class="loading-message">Loading...</p>}

      {nextPageToken.value && !loading.value && (
        <button
          class="btn btn-secondary load-more"
          onClick={() => loadPosts(nextPageToken.value)}
        >
          Load more
        </button>
      )}
    </div>
  );
}
