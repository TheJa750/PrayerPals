// frontend/src/lib/validation.js
// Validation functions for user input

// Validate password strength
export function validatePassword(password) {
    const errors = [];
    const minLength = 8;
    const hasUpperCase = /[A-Z]/.test(password);
    const hasLowerCase = /[a-z]/.test(password);
    const hasNumber = /[0-9]/.test(password);
    const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password);

    if (password.length < minLength) {
        errors.push(`Password must be at least ${minLength} characters long.`);
    }
    if (!hasUpperCase) {
        errors.push("Password must contain at least one uppercase letter.");
    }
    if (!hasLowerCase) {
        errors.push("Password must contain at least one lowercase letter.");
    }
    if (!hasNumber) {
        errors.push("Password must contain at least one number.");
    }
    if (!hasSpecialChar) {
        errors.push("Password must contain at least one special character.");
    }

    return {
        isValid: errors.length === 0,
        errors: errors
    };
}

// Validate username
export function validateUsername(username) {
    const errors = [];
    const minLength = 2;
    const maxLength = 25;
    const validChars = /^[a-zA-Z0-9_]+$/;

    if (username.length < minLength || username.length > maxLength) {
        errors.push(`Username must be between ${minLength} and ${maxLength} characters long.`);
    }
    if (!validChars.test(username)) {
        errors.push("Username can only contain letters, numbers, and underscores.");
    }

    return {
        isValid: errors.length === 0,
        errors: errors
    };
}

// Validate email format
export function validateEmail(email) {
    const errors = [];
    const emailRegex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;

    if (!email) {
        errors.push("Email is required");
        return {
            isValid: false,
            errors: errors
        };
    }

    if (email.length > 254) {
        errors.push("Email address is too long");
    }

    if (!emailRegex.test(email)) {
        errors.push("Please enter a valid email address");
    }

    const commonDomainTypos = {
        'gmial.com': 'gmail.com',
        'gmai.com': 'gmail.com',
        'yahooo.com': 'yahoo.com',
        'hotmial.com': 'hotmail.com'
    };

    const domain = email.split('@')[1];
    if (domain && commonDomainTypos[domain.toLowerCase()]) {
        errors.push(`Did you mean ${email.replace(domain, commonDomainTypos[domain.toLowerCase()])}?`);
    }

    return {
        isValid: errors.length === 0,
        errors: errors
    };
}

// Validate Group Invite Code
export function validateInviteCode(code) {
    const errors = [];
    const minLength = 1;
    const maxLength = 6;

    const validChars = /^[a-zA-Z0-9]+$/;

    if (code.length < minLength) {
        errors.push(`Invite code must be at least ${minLength} character long.`);
    }
    if (code.length > maxLength) {
        errors.push(`Invite code must be at most ${maxLength} characters long.`);
    }
    if (!validChars.test(code)) {
        errors.push("Invite code can only contain letters and numbers.");
    }

    return {
        isValid: errors.length === 0,
        errors: errors
    };
}