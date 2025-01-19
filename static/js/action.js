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
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ selectedId: selectedId })
    })
        .then(response => response.json())
        .then(data => {
            if (data.finished) {
                window.location.href = '/result';
            } else {
                updateBattle(data);
            }
        });
}

function updateBattle(data) {
    document.querySelector('.round-info').textContent =
        `${data.currentRound} (${data.matchNumber}/${data.totalMatches})`;

    const candidates = document.querySelectorAll('.candidate');
    candidates[0].onclick = () => selectCandidate(data.currentBattle.Candidate1.ID);
    candidates[0].querySelector('img').src = data.currentBattle.Candidate1.Image;
    candidates[0].querySelector('.name').textContent = data.currentBattle.Candidate1.Name;

    candidates[1].onclick = () => selectCandidate(data.currentBattle.Candidate2.ID);
    candidates[1].querySelector('img').src = data.currentBattle.Candidate2.Image;
    candidates[1].querySelector('.name').textContent = data.currentBattle.Candidate2.Name;
} 