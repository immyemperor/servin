// Code copy functionality
document.addEventListener('DOMContentLoaded', function() {
    // Add copy buttons to code blocks
    const codeBlocks = document.querySelectorAll('pre');
    
    codeBlocks.forEach(function(codeBlock) {
        // Skip if already has a copy button
        if (codeBlock.querySelector('.copy-button')) {
            return;
        }
        
        // Create copy button
        const copyButton = document.createElement('button');
        copyButton.className = 'copy-button';
        copyButton.innerHTML = `
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            <span class="copy-text">Copy</span>
        `;
        
        // Position the copy button
        codeBlock.style.position = 'relative';
        copyButton.style.position = 'absolute';
        copyButton.style.top = '0.75rem';
        copyButton.style.right = '0.75rem';
        copyButton.style.background = 'rgba(0, 0, 0, 0.7)';
        copyButton.style.color = 'white';
        copyButton.style.border = 'none';
        copyButton.style.borderRadius = '6px';
        copyButton.style.padding = '0.5rem';
        copyButton.style.fontSize = '0.75rem';
        copyButton.style.cursor = 'pointer';
        copyButton.style.display = 'flex';
        copyButton.style.alignItems = 'center';
        copyButton.style.gap = '0.25rem';
        copyButton.style.opacity = '0';
        copyButton.style.transition = 'opacity 0.2s ease';
        copyButton.style.zIndex = '10';
        
        // Show/hide copy button on hover
        codeBlock.addEventListener('mouseenter', function() {
            copyButton.style.opacity = '1';
        });
        
        codeBlock.addEventListener('mouseleave', function() {
            copyButton.style.opacity = '0';
        });
        
        // Copy functionality
        copyButton.addEventListener('click', async function() {
            const code = codeBlock.querySelector('code');
            const text = code ? code.textContent : codeBlock.textContent;
            
            try {
                await navigator.clipboard.writeText(text);
                
                // Show success feedback
                const originalContent = copyButton.innerHTML;
                copyButton.innerHTML = `
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M20 6L9 17l-5-5"></path>
                    </svg>
                    <span class="copy-text">Copied!</span>
                `;
                copyButton.style.background = 'rgba(16, 185, 129, 0.8)';
                
                setTimeout(function() {
                    copyButton.innerHTML = originalContent;
                    copyButton.style.background = 'rgba(0, 0, 0, 0.7)';
                }, 2000);
                
            } catch (err) {
                console.error('Failed to copy code:', err);
                
                // Fallback for older browsers
                const textArea = document.createElement('textarea');
                textArea.value = text;
                document.body.appendChild(textArea);
                textArea.focus();
                textArea.select();
                
                try {
                    document.execCommand('copy');
                    copyButton.querySelector('.copy-text').textContent = 'Copied!';
                    setTimeout(function() {
                        copyButton.querySelector('.copy-text').textContent = 'Copy';
                    }, 2000);
                } catch (fallbackErr) {
                    console.error('Fallback copy failed:', fallbackErr);
                }
                
                document.body.removeChild(textArea);
            }
        });
        
        codeBlock.appendChild(copyButton);
    });
});

// Search functionality for documentation
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('search-input');
    const navLinks = document.querySelectorAll('.nav-link');
    
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            const searchTerm = this.value.toLowerCase();
            
            navLinks.forEach(function(link) {
                const text = link.textContent.toLowerCase();
                const listItem = link.closest('.nav-section');
                
                if (text.includes(searchTerm) || searchTerm === '') {
                    link.style.display = 'block';
                    if (listItem) {
                        listItem.style.display = 'block';
                    }
                } else {
                    link.style.display = 'none';
                }
            });
            
            // Hide/show nav sections based on visible links
            document.querySelectorAll('.nav-section').forEach(function(section) {
                const visibleLinks = section.querySelectorAll('.nav-link[style="display: block"], .nav-link:not([style])');
                if (visibleLinks.length === 0 && searchTerm !== '') {
                    section.style.display = 'none';
                } else {
                    section.style.display = 'block';
                }
            });
        });
    }
});
