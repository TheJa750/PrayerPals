<script>
  import Banner from "./components/Banner.svelte";
  import BannerBtn from "./components/BannerBtn.svelte";
  import Home from "./pages/Home.svelte";
  import Login from "./pages/Login.svelte";
  import Signup from "./pages/Signup.svelte";
  import User from "./pages/User.svelte";
  import Group from "./pages/Group.svelte";
  import Post from "./pages/Post.svelte";
  import Logout from "./components/Logout.svelte";
  import GroupSettings from "./pages/GroupSettings.svelte";

  // Page state variable: 'home' | 'login' | 'signup'
  let page = "home";
  let groupId = null; // For group navigation
  let postId = null; // For post navigation

  // Navigation helper—pass to child components (for SPA-feeling in-app nav)
  function navigate(where, gId = null, pId = null) {
    page = where;
    if (gId !== null) {
      groupId = gId;
    }
    if (pId !== null) {
      postId = pId;
    }
  }
</script>

<main>
  {#if page === "home"}
    <Banner>
      <BannerBtn page="login" text="Log In" slot="slot5" {navigate}></BannerBtn>
      <BannerBtn page="signup" text="Sign Up" slot="slot4" {navigate}
      ></BannerBtn>
    </Banner>
  {:else if page === "login"}
    <Banner>
      <BannerBtn page="signup" text="Sign Up" slot="slot4" {navigate}
      ></BannerBtn>
      <BannerBtn page="home" text="Back to Home" slot="slot5" {navigate}
      ></BannerBtn>
    </Banner>
  {:else if page === "signup"}
    <Banner>
      <BannerBtn page="login" text="Log In" slot="slot4" {navigate}></BannerBtn>
      <BannerBtn page="home" text="Back to Home" slot="slot5" {navigate}
      ></BannerBtn>
    </Banner>
  {:else if page === "user"}
    <Banner>
      <Logout {navigate} slot="slot5" />
      <h2 slot="slot1" class="welcome-message">Welcome Back!</h2>
    </Banner>
  {:else}
    <Banner>
      <Logout {navigate} slot="slot5" />
    </Banner>
  {/if}
  <div class="container">
    {#if page === "home"}
      <Home />
    {:else if page === "login"}
      <Login {navigate} />
    {:else if page === "signup"}
      <Signup {navigate} />
    {:else if page === "user"}
      <User {navigate} />
    {:else if page === "group"}
      <Group {navigate} {groupId} />
    {:else if page === "post"}
      <Post {navigate} {groupId} {postId} />
    {:else if page === "groupSettings"}
      <GroupSettings {navigate} {groupId} />
    {:else}
      <p>Page not found</p>
    {/if}
  </div>
</main>
