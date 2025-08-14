<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import { REFRESH_ERROR_MESSAGE } from "../lib/api";
    import CreatePostModal from "../components/CreatePostModal.svelte";
    import DeletePostModal from "../components/DeletePostModal.svelte";
    import GroupMembersModal from "../components/GroupMembersModal.svelte";

    export let navigate;
    export let groupId;

    // Group data
    let group = null;
    let groupPosts = [];
    let isLoadingGroup = true;
    let isLoadingPosts = true;
    let loadGroupError = "";
    let loadPostError = "";
    let userRole = "member";
    let isAdmin = false;
    let userId = localStorage.getItem("userId");

    // Create post modal data
    let newPostContent = "";
    let isCreatingPost = false;
    let createPostError = "";

    // Delete post modal data
    let deletePostId = null;
    let isDeletingPost = false;
    let deletePostError = "";

    // Members modal data
    let members = [];
    let isLoadingMembers = false;
    let loadMembersError = "";

    // Member moderation modal data
    let targetMember = null;
    let action = "";
    let reason = "";
    let moderationError = "";
    let deletePosts = false;

    //Modal states
    let showCreatePostModal = false;
    let showDeletePostModal = false;
    let showMembersModal = false;
    let showErrorModal = false;
    let showMemberModerationModal = false;

    // Load group data on mount
    onMount(async () => {
        await checkUserRole();
        loadGroupData();
        loadGroupPosts();
    });

    // Data Loading Functions
    async function loadGroupData() {
        try {
            isLoadingGroup = true;
            loadGroupError = "";

            const groupData = await apiRequest(`/groups/${groupId}`, "GET");
            group = groupData;
        } catch (error) {
            console.error("Error loading group data:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            loadGroupError =
                "Failed to load group data. Please try again later.";
        } finally {
            isLoadingGroup = false;
        }
    }

    async function loadGroupPosts() {
        try {
            isLoadingPosts = true;
            const posts = await apiRequest(`/groups/${groupId}/posts`, "GET");
            groupPosts = posts || [];
        } catch (error) {
            console.error("Error loading group posts:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            loadPostError =
                "Failed to load prayer requests. Please try again later.";
        } finally {
            isLoadingPosts = false;
        }
    }

    async function checkUserRole() {
        try {
            const roleData = await apiRequest(
                `/groups/${groupId}/members/${userId}`,
                "GET",
            );

            userRole = roleData.role || "member";
            isAdmin = userRole === "admin";
        } catch (error) {
            console.error("Error checking user role:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            showErrorModal = true;
            loadGroupError = error.message || "Failed to check user role.";

            setTimeout(() => {
                navigate("user");
            }, 2500); // Delay to show error message before returning

            return;
        }
    }

    // Modal Functions
    function openCreatePostModal() {
        showCreatePostModal = true;
        newPostContent = "";
        createPostError = "";
    }

    function closeCreatePostModal() {
        showCreatePostModal = false;
    }

    function openDeletePostModal(id) {
        deletePostId = id;
        isDeletingPost = false;
        showDeletePostModal = true;
        deletePostError = "";
    }

    function closeDeletePostModal() {
        showDeletePostModal = false;
        deletePostId = null;
        isDeletingPost = false;
    }

    function openMembersModal() {
        showMembersModal = true;
        fetchMembers();
    }

    function closeMembersModal() {
        showMembersModal = false;
    }

    function handleModalContentUpdate(event) {
        newPostContent = event.detail;
    }

    async function handleCreatePost(event) {
        event.preventDefault();
        if (!newPostContent.trim()) {
            createPostError = "Post content cannot be empty.";
            return;
        }

        isCreatingPost = true;
        createPostError = "";

        try {
            const newPost = await apiRequest(
                `/groups/${groupId}/posts`,
                "POST",
                { content: newPostContent },
            );

            closeCreatePostModal();
            await loadGroupPosts(); // Reload posts after creating a new one
        } catch (error) {
            console.error("Error creating post:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            createPostError =
                error.message || "Failed to create post. Please try again.";
        } finally {
            isCreatingPost = false;
        }
    }

    async function handleDeletePost(event) {
        event.preventDefault();

        if (!deletePostId) {
            deletePostError = "No post selected for deletion.";
            return;
        }

        isDeletingPost = true;
        deletePostError = "";

        try {
            await apiRequest(
                `/groups/${groupId}/posts/${deletePostId}`,
                "DELETE",
            );

            closeDeletePostModal();
            await loadGroupPosts(); // Reload posts after deletion
        } catch (error) {
            console.error("Error deleting post:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            deletePostError =
                error.message || "Failed to delete post. Please try again.";
        } finally {
            isDeletingPost = false;
        }
    }

    function formatTimestamp(timestamp) {
        const date = new Date(timestamp);
        return (
            date.toLocaleDateString() +
            " at " +
            date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
        );
    }

    async function fetchMembers() {
        isLoadingMembers = true;
        loadMembersError = "";

        try {
            const membersData = await apiRequest(
                `/groups/${groupId}/members`,
                "GET",
            );
            members = membersData || [];
        } catch (error) {
            console.error("Error loading group members:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            loadMembersError =
                "Failed to load group members. Please try again later.";
        } finally {
            isLoadingMembers = false;
        }
    }

    function showModerateMember(event) {
        event.preventDefault();
        targetMember = event.detail;

        // Close member modal and open confirmation modal
        showMembersModal = false;
        showMemberModerationModal = true;
    }

    function closeMemberModerationModal() {
        showMemberModerationModal = false;
        targetMember = null;
        action = "";
        reason = "";
        moderationError = "";
    }

    async function handleModerationAction(event) {
        event.preventDefault();

        if (!action || !targetMember) {
            return; // Ensure action and target member are selected
        }

        try {
            // Make moderation API call
            const reqData = {
                action: action,
                reason: reason || "No reason provided",
            };
            await apiRequest(
                `/groups/${groupId}/members/${targetMember.user_id}/moderate`,
                "PUT",
                reqData,
            );

            // Remove posts if checkbox is checked
            if (deletePosts) {
                await apiRequest(
                    `/groups/${groupId}/members/${targetMember.user_id}/remove-content`,
                    "PUT",
                );
            }

            await loadGroupPosts(); // Reload posts after moderation
            closeMemberModerationModal();
            await fetchMembers(); // Refresh members list after moderation
            openMembersModal(); // Reopen members modal to show updated list
        } catch (error) {
            console.error("Error moderating member:", error);
            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }
            moderationError =
                error.message || "Failed to moderate member. Please try again.";
        }
    }
</script>

<div class="group-container">
    <!-- Group Header -->
    <div class="group-header">
        <div class="header-row">
            <button class="back-button" on:click={() => navigate("user")}>
                ←
            </button>
            <div class="header-content">
                {#if isLoadingGroup}
                    <h1>Loading group...</h1>
                {:else if loadGroupError}
                    <h1>Error: {loadGroupError}</h1>
                {:else if group}
                    <h1>{group.name}</h1>
                    {#if group.description}
                        <p>{group.description}</p>
                    {/if}
                {:else}
                    <h1>Group not found</h1>
                {/if}
            </div>
            <div class="header-spacer"></div>
        </div>
    </div>

    <!-- Main Grid Layout -->
    <div class="group-grid">
        <!-- Left: Group Actions -->
        <section class="group-section side-section">
            <h2>Group Actions</h2>
            <div class="actions-grid">
                <button class="action-button" on:click={openCreatePostModal}
                    >New Request</button
                >
                <button class="action-button" on:click={openMembersModal}
                    >View Members</button
                >
                <button class="action-button">Leave Group</button>
            </div>
        </section>

        <!-- Center: Prayer Request Feed -->
        <section class="group-section main-section">
            <h2>Prayer Request</h2>

            {#if isLoadingPosts}
                <p>Loading prayer requests...</p>
            {:else if groupPosts.length === 0}
                <p>Be the first to post your prayer request!</p>
            {:else}
                {#each groupPosts as post}
                    <div
                        class="post-card clickable"
                        on:click={() => navigate("post", groupId, post.id)}
                        on:keydown={(e) =>
                            e.key === "Enter" &&
                            navigate("post", groupId, post.id)}
                        tabindex="0"
                        role="button"
                        aria-label="View post by {post.username}"
                    >
                        <div class="post-header">
                            <h3>{post.author}</h3>
                            {#if isAdmin || post.user_id === userId}
                                <button
                                    class="close-button"
                                    on:click|stopPropagation={() =>
                                        openDeletePostModal(post.id)}
                                >
                                    ×
                                </button>
                            {/if}
                        </div>
                        <div class="post-body">
                            <p>{post.content}</p>
                        </div>
                        <div class="post-footer">
                            <span class="post-timestamp"
                                >{formatTimestamp(post.created_at)}</span
                            >
                            <span class="post-comments"
                                >{post.comment_count || 0}</span
                            >
                        </div>
                    </div>
                {/each}
            {/if}
        </section>

        <!-- Right: Group Info -->
        <section class="group-section side-section">
            <h2>Group Info</h2>
            <p>Ground rules and announcements go here</p>
        </section>
    </div>
</div>

<!-- Error Modal -->
{#if showErrorModal}
    <div class="modal-overlay" role="dialog" aria-modal="true" tabindex="-1">
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
        <div
            class="modal-content text-center h-50 w-50"
            on:click|stopPropagation
            role="document"
        >
            <h2>Error:</h2>
            <p>{loadGroupError}</p>
            <p>Returning to your dashboard...</p>
        </div>
    </div>
{/if}

<!-- Create Post Modal -->
{#if showCreatePostModal}
    <CreatePostModal
        {newPostContent}
        isCreating={isCreatingPost}
        error={createPostError}
        on:updateContent={handleModalContentUpdate}
        on:close={closeCreatePostModal}
        on:submit={handleCreatePost}
    />
{/if}

<!-- Delete Post Modal -->
{#if showDeletePostModal}
    <DeletePostModal
        isDeleting={isDeletingPost}
        error={deletePostError}
        on:close={closeDeletePostModal}
        on:submit={handleDeletePost}
    />
{/if}

<!-- Group Members Modal -->
{#if showMembersModal}
    <GroupMembersModal
        {members}
        inviteCode={group ? group.invite_code : ""}
        isLoading={isLoadingMembers}
        error={loadMembersError}
        {isAdmin}
        on:close={closeMembersModal}
        on:remove={showModerateMember}
    />
{/if}

<!-- Member Moderation Modal -->
{#if showMemberModerationModal}
    <div class="modal-overlay" role="dialog" aria-modal="true" tabindex="-1">
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
        <div
            class="modal-content text-center h-50 w-50"
            on:click|stopPropagation
            role="document"
        >
            <div class="modal-header">
                <h2>Moderate {targetMember.username}</h2>
            </div>
            <div class="modal-body">
                <form>
                    <div class="form-row">
                        <label for="moderation-action">Select Action:</label>
                        <select
                            id="moderation-action"
                            bind:value={action}
                            required
                        >
                            <option value="" disabled selected
                                >Select action</option
                            >
                            <option value="kick"
                                >Kick from Group (lasts 1 week)</option
                            >
                            <option value="ban"
                                >Ban from Group (permanent)</option
                            >
                        </select>
                    </div>
                    <div class="form-row">
                        <label for="moderation-reason">Reason:</label>
                        <textarea
                            id="moderation-reason"
                            rows="3"
                            placeholder="Enter reason for moderation action"
                            required
                            bind:value={reason}
                        ></textarea>
                    </div>
                    <div class="form-row">
                        <label for="remove-posts"
                            >Remove all posts by {targetMember.username}</label
                        >
                        <input
                            type="checkbox"
                            id="remove-posts"
                            bind:checked={deletePosts}
                        />
                    </div>
                    <div class="modal-actions">
                        <button
                            type="button"
                            on:click={closeMemberModerationModal}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            on:click={handleModerationAction}
                            disabled={!action || !reason}
                        >
                            Confirm
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
{/if}
