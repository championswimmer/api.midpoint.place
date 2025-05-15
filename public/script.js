// Smooth scroll for navigation links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Add intersection observer for fade-in animations
const observerOptions = {
    root: null,
    rootMargin: '0px',
    threshold: 0.1
};

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.classList.add('fade-in');
            observer.unobserve(entry.target);
        }
    });
}, observerOptions);

// Observe all feature cards and step cards
document.querySelectorAll('.feature-card, .step-card').forEach(card => {
    observer.observe(card);
});

// Handle Get Started button clicks
document.querySelectorAll('button').forEach(button => {
    button.addEventListener('click', () => {
        // This is where you would handle the sign-up flow
        // For now, we'll just add a simple animation
        button.classList.add('clicked');
        setTimeout(() => {
            button.classList.remove('clicked');
        }, 200);
    });
}); 