import { Head } from "fresh/runtime";
import { define } from "../utils.ts";
import PostManager from "../islands/PostManager.tsx";

export default define.page(function Home() {
  return (
    <div class="container">
      <Head>
        <title>Guestbook</title>
      </Head>
      <header class="header">
        <h1>Guestbook</h1>
      </header>
      <main>
        <PostManager />
      </main>
    </div>
  );
});
