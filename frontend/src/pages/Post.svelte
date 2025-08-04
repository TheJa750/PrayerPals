<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import { REFRESH_ERROR_MESSAGE } from "../lib/api";

    export let navigate;
    export let groupId;
    export let postId;

    // Post and comments data
    let post = null;
    let comments = [];
    let isLoadingPost = true;
    let loadPostError = ""; // Changed from errorMessage

    // Add comment form
    let newCommentContent = "";
    let isAddingComment = false;
    let commentError = "";

    // Timestamp formatter (reuse from Group.svelte)
    function formatTimestamp(timestamp) {
        const date = new Date(timestamp);
        return (
            date.toLocaleDateString() +
            " at " +
            date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
        );
    }

    onMount(() => {
        loadPost();
    });

    async function loadPost() {
        try {
            isLoadingPost = true;
            loadPostError = "";

            const postData = await apiRequest(
                `/groups/${groupId}/posts/${postId}/comments`,
                "GET",
            );
            post = postData;
            comments = post.comments || []; // Ensure comments is always an array
        } catch (error) {
            console.error("Error loading post:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            loadPostError = "Failed to load post";
        } finally {
            isLoadingPost = false;
        }
    }

    async function addComment() {
        commentError = "";
        if (!newCommentContent.trim()) return;

        isAddingComment = true;
        try {
            const response = await apiRequest(
                `/groups/${groupId}/posts/${postId}/comments`,
                "POST",
                {
                    content: newCommentContent.trim(),
                },
            );

            newCommentContent = ""; // Clear the input field
            await loadPost(); // Reload the post to include the new comment
        } catch (error) {
            console.error("Error adding comment:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            commentError = error.message || "Failed to add comment";
        } finally {
            isAddingComment = false;
        }
    }
</script>

<div class="post-container">
    <!-- Back button and header -->
    <div class="post-header">
        <button class="back-button" on:click={() => navigate("group", groupId)}>
            ←
        </button>
        <h1>Prayer Request</h1>
    </div>

    <!-- Main post -->
    {#if isLoadingPost}
        <div class="post-card">
            <p>Loading post...</p>
        </div>
    {:else if loadPostError}
        <div class="post-card">
            <p class="error-message">{loadPostError}</p>
        </div>
    {:else if post}
        <div class="post-card main-post">
            <div class="post-author">
                <h2>{post.author}</h2>
            </div>
            <div class="post-content">
                <p>{post.content}</p>
            </div>
            <div class="post-timestamp">
                <span>{formatTimestamp(post.created_at)}</span>
            </div>
        </div>
    {/if}

    <!-- Comments section -->
    <div class="comments-section">
        <h3>Comments ({comments.length})</h3>

        {#if comments.length === 0}
            <p>No comments yet. Be the first to comment!</p>
        {:else}
            {#each comments as comment}
                <div class="comment-card">
                    <div class="post-author">
                        <strong>{comment.author}</strong>
                    </div>
                    <div class="post-content">
                        <p>{comment.content}</p>
                    </div>
                    <div class="post-timestamp">
                        <small>{formatTimestamp(comment.created_at)}</small>
                    </div>
                </div>
            {/each}
        {/if}
    </div>

    <!-- Add comment form -->
    <div class="add-comment-section">
        <h3>Add Comment</h3>
        <form on:submit|preventDefault={addComment}>
            <textarea
                bind:value={newCommentContent}
                placeholder="Write your comment here..."
                rows="5"
                required
            ></textarea>
            <button
                type="submit"
                disabled={isAddingComment || newCommentContent.trim() === ""}
                >{isAddingComment ? "Adding..." : "Add Comment"}
            </button>
            {#if commentError}
                <p class="error-message">{commentError}</p>
            {/if}
        </form>
    </div>
</div>
