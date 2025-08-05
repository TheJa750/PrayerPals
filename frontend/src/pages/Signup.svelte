<script>
    import { apiRequest } from "../lib/api";
    import {
        validateEmail,
        validatePassword,
        validateUsername,
    } from "../lib/validation";

    export let navigate;

    let username = "";
    let email = "";
    let password = "";
    let showPassword = false;

    let isSubmitting = false;
    let errorMessage = "";
    let successMessage = "";

    // Validation states
    let usernameValidation = { isValid: true, errors: [] };
    let emailValidation = { isValid: true, errors: [] };
    let passwordValidation = { isValid: true, errors: [] };

    // Reactive validation checks
    $: usernameValidation = username
        ? validateUsername(username)
        : { isValid: true, errors: [] };
    $: emailValidation = email
        ? validateEmail(email)
        : { isValid: true, errors: [] };
    $: passwordValidation = password
        ? validatePassword(password)
        : { isValid: true, errors: [] };

    // Check if form is valid
    $: isFormValid =
        usernameValidation.isValid &&
        emailValidation.isValid &&
        passwordValidation.isValid;

    async function handleSubmit(event) {
        event.preventDefault();

        // Validate form before submission
        if (!isFormValid) {
            errorMessage =
                "Please fix the errors in the form before submitting.";
            return;
        }

        isSubmitting = true;
        errorMessage = "";
        successMessage = "";

        try {
            // create account
            const userData = {
                username: username,
                email: email,
                password: password,
            };

            const response = await apiRequest("/users", "POST", userData);

            successMessage = "Account created! Logging you in...";

            // Setup the data before clearing the form
            const loginData = {
                email: email,
                password: password,
            };

            try {
                const loginResponse = await apiRequest(
                    "/login",
                    "POST",
                    loginData,
                );

                //Successfully logged in, navigate to home
                successMessage = "Logged in! Redirecting to your dashboard...";

                // clear form fields
                username = "";
                email = "";
                password = "";

                // Navigate to dashboard after brief delay
                setTimeout(() => {
                    navigate("user");
                }, 1000);
            } catch (loginError) {
                successMessage = "";
                errorMessage =
                    "Account created, but failed to log in. Please log in manually.";
            }
        } catch (signupError) {
            errorMessage =
                "Failed to create account: " +
                (signupError.message || "Unknown error");
        } finally {
            isSubmitting = false;
        }
    }
</script>

<h2 class="text-center mb-1">Create Account</h2>
<form class:form-submitting={isSubmitting} on:submit={handleSubmit}>
    <div class="form-row">
        <label for="signup-username">Username:</label>
        <input
            type="text"
            id="signup-username"
            required
            bind:value={username}
        />
    </div>
    {#if !usernameValidation.isValid && username}
        <div class="form-row validation-container">
            {#each usernameValidation.errors as error}
                <p class="validation-error">{error}</p>
            {/each}
        </div>
    {/if}

    <div class="form-row">
        <label for="signup-email">Email:</label>
        <input type="email" id="signup-email" required bind:value={email} />
    </div>
    {#if !emailValidation.isValid && email}
        <div class="form-row validation-container">
            {#each emailValidation.errors as error}
                <p class="validation-error">{error}</p>
            {/each}
        </div>
    {/if}

    <div class="form-row">
        <label for="signup-password">Password:</label>
        <input
            type={showPassword ? "text" : "password"}
            id="signup-password"
            required
            bind:value={password}
        />
    </div>
    {#if !passwordValidation.isValid && password}
        <div class="validation-container">
            {#each passwordValidation.errors as error}
                <p class="validation-error">{error}</p>
            {/each}
        </div>
    {/if}

    {#if errorMessage || successMessage}
        <div class="form-row">
            {#if successMessage}
                <div class="success-message text-center">{successMessage}</div>
            {/if}
            {#if errorMessage}
                <div class="error-message text-center">{errorMessage}</div>
            {/if}
        </div>
    {/if}
    <div class="form-row">
        <button
            type="button"
            class="toggle-password"
            on:click={() => (showPassword = !showPassword)}
            aria-label={showPassword ? "Hide password" : "Show password"}
        >
            {showPassword ? "🙈 Hide" : "👁️ Show"}
        </button>
        <button disabled={isSubmitting}
            >{isSubmitting ? "Creating Account..." : "Submit"}</button
        >
    </div>
</form>
