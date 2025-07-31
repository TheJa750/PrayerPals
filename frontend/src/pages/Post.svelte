<script>
    export let postId;
    export let groupId;
    export let navigate;

    let comments = [];
    let loading = true;
    let error = "";

    // Fetch comments on component mount or when postId changes
    $: if (postId && groupId) {
        loadComments();
    }

    async function loadComments() {
        loading = true;
        error = "";
        comments = [];
        try {
            const res = await fetch(
                `http://localhost:8080/api/groups/${encodeURIComponent(groupId)}/posts/${encodeURIComponent(postId)}/comments`,
                { credentials: "include" },
            );
            if (res.ok) {
                comments = await res.json();
            } else {
                error = "Failed to load comments: " + (await res.text());
            }
        } catch (err) {
            error = "Error: " + err.message;
        }
        loading = false;
    }
</script>

<main>
    <div class="container">
        <h1>Post Discussion</h1>

        {#if loading}
            <p>Loading comments...</p>
        {:else if error}
            <p style="color:red;">{error}</p>
        {:else if comments.length === 0}
            <p>No comments yet. Be the first to leave encouragement!</p>
        {:else}
            <div id="commentsContainer">
                {#each comments as comment}
                    <div class="group-card">
                        <strong>{comment.user_id}</strong><br />
                        <span>{comment.content}</span>
                        <div style="font-size:0.8em; color:#999;">
                            {comment.created_at}
                        </div>
                    </div>
                {/each}
            </div>
        {/if}

        <div style="text-align:center; margin-top: 1em;">
            <button on:click={() => navigate("group", groupId)}
                >Back to Group</button
            >
        </div>
    </div>
</main>
