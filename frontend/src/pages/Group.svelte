<script>
    export let navigate;
    export let groupId;

    // Stub data for posts (replace with real API call in onMount)
    let posts = [
        {
            id: "post1",
            author: "Alice",
            content: "Pray for my exams next week!",
        },
        { id: "post2", author: "Bob", content: "Health for my family." },
    ];

    function leaveGroup() {
        // Call API to leave group, then go back to dashboard
        navigate("dashboard");
    }

    function makePost() {
        navigate("new-post", groupId);
    }

    function viewMembers() {
        navigate("members", groupId);
    }
</script>

<main>
    <div class="container">
        <h1>Group Dashboard</h1>
        <div id="postsContainer">
            {#if posts.length === 0}
                <p>No posts yet. Be the first to share a prayer request!</p>
            {:else}
                {#each posts as post}
                    <div class="group-card">
                        <strong>{post.author}</strong><br />
                        <span>{post.content}</span>
                    </div>
                {/each}
            {/if}
        </div>
        <div style="text-align:center;">
            <button on:click={makePost}>Make Post</button>
            <button on:click={viewMembers}>Members</button>
            <button on:click={leaveGroup}>Leave Group</button>
            <button on:click={() => navigate("dashboard")}
                >Back to Dashboard</button
            >
        </div>
    </div>
</main>
