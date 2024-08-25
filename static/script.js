document.addEventListener('DOMContentLoaded', function() {
    const form = document.querySelector('form');
    const textInput = document.getElementById('text');
    const bannerSelect = document.getElementById('banner');
    const previewArea = document.getElementById('preview');

    function updatePreview() {
        const text = textInput.value;
        const banner = bannerSelect.value;
        
        fetch('/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `text=${encodeURIComponent(text)}&banner=${encodeURIComponent(banner)}`
        })
        .then(response => response.text())
        .then(html => {
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            const output = doc.querySelector('pre');
            if (output) {
                previewArea.innerHTML = output.innerHTML;
            }
        });
    }

    textInput.addEventListener('input', updatePreview);
    bannerSelect.addEventListener('change', updatePreview);
});
