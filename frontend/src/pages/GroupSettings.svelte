<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import { REFRESH_ERROR_MESSAGE } from "../lib/api";

    export let navigate;
    export let groupId;

    // Group data variables
    let group = null;
    let isLoadingGroup = false;
    let loadGroupError = "";

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

    onMount(() => {
        if (groupId) {
            loadGroupData();
        } else {
            navigate("user");
        }
    });
</script>

<div class="group-settings">
    <div class="header-row">
        <button class="back-button" on:click={() => navigate("group", groupId)}>
            ←
        </button>
        <div class="header-content">
            <h2>Group Settings</h2>
            <p>Manage your group settings here.</p>
        </div>
    </div>
    {#if isLoadingGroup}
        <p>Loading group settings...</p>
    {:else if loadGroupError}
        <p class="error">{loadGroupError}</p>
    {:else if group}
        <div class="settings-card">
            <h2>{group.name}</h2>
            <p>{group.description}</p>
            <p>{group.invite_code}</p>
        </div>
    {/if}
</div>
