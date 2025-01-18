function createHeart(event) {
    const hearts = ['❤️', '💖', '💝', '💕', '💗'];
    const numHearts = 5;
    const container = event.currentTarget;
    const rect = container.getBoundingClientRect();

    for (let i = 0; i < numHearts; i++) {
        setTimeout(() => {
            const heart = document.createElement('div');
            heart.className = 'heart';
            heart.innerHTML = hearts[Math.floor(Math.random() * hearts.length)];

            const x = event.clientX - rect.left + (Math.random() - 0.5) * 60;
            const y = event.clientY - rect.top + (Math.random() - 0.5) * 60;

            heart.style.cssText = `left: ${x}px; top: ${y}px;`;

            container.appendChild(heart);

            heart.addEventListener('animationend', () => {
                heart.remove();
            });
        }, i * 100);
    }
}

function selectCandidate(id) {
    const selectedId = parseInt(id, 10);
    fetch('/game/select', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ selectedId })
    })
        .then(response => response.json())
        .then(data => {
            if (data.finished) {
                window.location.href = '/result';
            } else {
                window.location.reload();
            }
        });
} 