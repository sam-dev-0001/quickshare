// Setup custom key bindings and utility listeners
document.addEventListener('DOMContentLoaded', () => {
    initCharacterCounter();
    initTabOverride();
    initKeyboardShortcuts();
    initTimestamps();
});

// Toast notification helper
function showToast(message, type = 'success') {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    
    // Icon mapping
    let icon = '✨';
    if (type === 'error') icon = '❌';
    else if (type === 'success') icon = '✓';

    toast.innerHTML = `<span class="toast-icon">${icon}</span> <span class="toast-message">${message}</span>`;
    container.appendChild(toast);

    // Auto delete after 3 seconds
    setTimeout(() => {
        toast.style.animation = 'slideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1) reverse forwards';
        setTimeout(() => {
            toast.remove();
        }, 300);
    }, 3000);
}

// Clipboard Copy logic
function copySnippetToClipboard() {
    const codeElement = document.getElementById('snippet-code-block');
    if (!codeElement) {
        showToast('Failed to find code element to copy', 'error');
        return;
    }

    const textToCopy = codeElement.innerText;
    
    navigator.clipboard.writeText(textToCopy)
        .then(() => {
            showToast('Snippet copied to clipboard!');
        })
        .catch(err => {
            console.error('Failed to copy text: ', err);
            showToast('Failed to copy. Please copy manually.', 'error');
        });
}

// Character counter inside text editor
function initCharacterCounter() {
    const textarea = document.getElementById('snippet-textarea');
    const counter = document.getElementById('char-count');
    
    if (!textarea || !counter) return;

    const updateCount = () => {
        const count = textarea.value.length;
        counter.textContent = count.toLocaleString();
    };

    textarea.addEventListener('input', updateCount);
    updateCount(); // Initial load
}

// Support TAB indentation inside the text editor
function initTabOverride() {
    const textarea = document.getElementById('snippet-textarea');
    if (!textarea) return;

    textarea.addEventListener('keydown', (e) => {
        if (e.key === 'Tab') {
            e.preventDefault();
            const start = textarea.selectionStart;
            const end = textarea.selectionEnd;

            // Set textarea value to: text before caret + 4 spaces + text after caret
            textarea.value = textarea.value.substring(0, start) + "    " + textarea.value.substring(end);

            // Put caret in correct position
            textarea.selectionStart = textarea.selectionEnd = start + 4;
        }
    });
}

// Global keyboard shortcuts (Ctrl + S or Cmd + S to trigger share)
function initKeyboardShortcuts() {
    document.addEventListener('keydown', (e) => {
        if ((e.ctrlKey || e.metaKey) && e.key === 's') {
            const form = document.getElementById('snippet-form');
            if (form) {
                e.preventDefault();
                // Trigger form submission
                form.dispatchEvent(new Event('submit', { cancelable: true, bubbles: true }));
            }
        }
    });
}

// Format raw ISO timestamps to localized human-readable format
function initTimestamps() {
    const timestamps = document.querySelectorAll('.timestamp');
    timestamps.forEach(el => {
        const rawTime = el.textContent;
        if (!rawTime) return;
        try {
            const date = new Date(rawTime);
            if (!isNaN(date.getTime())) {
                el.textContent = date.toLocaleString(undefined, {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                });
            }
        } catch (e) {
            console.error('Failed to parse timestamp:', e);
        }
    });
}
