// Wait for DOM to load
document.addEventListener('DOMContentLoaded', function() {
    // Smooth scrolling for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            
            const target = document.querySelector(this.getAttribute('href'));
            if (!target) return;
            
            window.scrollTo({
                top: target.offsetTop - 80, // Offset for fixed header
                behavior: 'smooth'
            });
        });
    });
    
    // Animation on scroll
    const animateOnScroll = function() {
        const items = document.querySelectorAll('.feature-card, .testimonial, .step');
        
        items.forEach(item => {
            const itemPosition = item.getBoundingClientRect().top;
            const screenPosition = window.innerHeight / 1.2;
            
            if (itemPosition < screenPosition) {
                item.classList.add('animate-in');
            }
        });
    };
    
    // Add CSS for animation classes
    const style = document.createElement('style');
    style.textContent = `
        .feature-card, .testimonial, .step {
            opacity: 0;
            transform: translateY(20px);
            transition: opacity 0.6s ease, transform 0.6s ease;
        }
        
        .animate-in {
            opacity: 1;
            transform: translateY(0);
        }
        
        .testimonial.animate-in:nth-child(2) {
            transform: rotate(1deg);
        }
        
        .testimonial.animate-in:nth-child(3) {
            transform: rotate(-1deg);
        }
    `;
    document.head.appendChild(style);
    
    // Run animation check on scroll
    window.addEventListener('scroll', animateOnScroll);
    
    // Initial check for items in viewport
    animateOnScroll();
    
    // Funky cursor effect for GenZ vibe
    const createCursorEffect = function() {
        const cursor = document.createElement('div');
        cursor.classList.add('cursor-effect');
        document.body.appendChild(cursor);
        
        const cursorStyle = document.createElement('style');
        cursorStyle.textContent = `
            .cursor-effect {
                width: 40px;
                height: 40px;
                border-radius: 50%;
                background: rgba(255, 94, 125, 0.2);
                position: fixed;
                transform: translate(-50%, -50%);
                pointer-events: none;
                z-index: 9999;
                transition: transform 0.1s, width 0.3s, height 0.3s, background 0.3s;
                transform-origin: 100% 100%;
            }
            
            .cursor-effect.active {
                transform: translate(-50%, -50%) scale(0.7);
                background: rgba(255, 94, 125, 0.5);
            }
            
            .btn-primary:hover ~ .cursor-effect,
            .btn-secondary:hover ~ .cursor-effect,
            a:hover ~ .cursor-effect {
                transform: translate(-50%, -50%) scale(1.5);
                background: rgba(255, 94, 125, 0.1);
            }
        `;
        document.head.appendChild(cursorStyle);
        
        document.addEventListener('mousemove', e => {
            cursor.style.left = e.clientX + 'px';
            cursor.style.top = e.clientY + 'px';
        });
        
        document.addEventListener('mousedown', () => {
            cursor.classList.add('active');
        });
        
        document.addEventListener('mouseup', () => {
            cursor.classList.remove('active');
        });
    };
    
    // Only enable cursor effect on non-touch devices
    if (!('ontouchstart' in window)) {
        createCursorEffect();
    }
    
    // Add interactive elements to show off the midpoint concept
    const createInteractiveDemo = function() {
        const demoContainer = document.querySelector('.md\\:w-1\\/2:last-child');
        
        // Only proceed if we found the container
        if (!demoContainer) return;
        
        // Create the interactive map placeholder
        const interactiveDemo = document.createElement('div');
        interactiveDemo.classList.add('interactive-demo');
        interactiveDemo.innerHTML = `
            <div class="demo-controls">
                <button class="demo-btn active" data-points="3">3 Friends</button>
                <button class="demo-btn" data-points="4">4 Friends</button>
                <button class="demo-btn" data-points="5">5 Friends</button>
            </div>
            <div class="demo-map">
                <div class="map-container">
                    <div class="map-point map-midpoint"></div>
                </div>
            </div>
        `;
        
        demoContainer.appendChild(interactiveDemo);
        
        // Add styles
        const demoStyle = document.createElement('style');
        demoStyle.textContent = `
            .interactive-demo {
                margin-top: 1rem;
                border-radius: 1rem;
                overflow: hidden;
                background: var(--color-card);
                border: 2px solid var(--color-border);
            }
            
            .demo-controls {
                display: flex;
                padding: 0.5rem;
                background: var(--color-card-hover);
                border-bottom: 2px solid var(--color-border);
            }
            
            .demo-btn {
                background: var(--color-card);
                border: 1px solid var(--color-border);
                border-radius: 0.5rem;
                padding: 0.5rem 1rem;
                margin-right: 0.5rem;
                cursor: pointer;
                font-weight: 500;
                transition: all 0.3s ease;
            }
            
            .demo-btn.active {
                background: var(--color-primary);
                color: white;
                border-color: var(--color-primary-dark);
            }
            
            .demo-map {
                height: 200px;
                position: relative;
                overflow: hidden;
            }
            
            .map-container {
                position: absolute;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
            }
            
            .map-point {
                position: absolute;
                width: 20px;
                height: 20px;
                border-radius: 50%;
                background: var(--color-secondary);
                transform: translate(-50%, -50%);
                z-index: 1;
                box-shadow: 0 0 0 3px rgba(51, 191, 255, 0.3);
            }
            
            .map-midpoint {
                background: var(--color-primary);
                width: 25px;
                height: 25px;
                box-shadow: 0 0 0 5px rgba(255, 94, 125, 0.3);
                z-index: 2;
            }
        `;
        document.head.appendChild(demoStyle);
        
        // Function to generate random positions
        const generatePoints = function(numPoints) {
            const mapContainer = document.querySelector('.map-container');
            // Clear existing points except midpoint
            document.querySelectorAll('.map-point:not(.map-midpoint)').forEach(el => el.remove());
            
            const width = mapContainer.clientWidth;
            const height = mapContainer.clientHeight;
            const midpoint = document.querySelector('.map-midpoint');
            
            // Generate random points
            let pointsData = [];
            for (let i = 0; i < numPoints; i++) {
                const x = Math.random() * (width - 60) + 30;
                const y = Math.random() * (height - 60) + 30;
                pointsData.push({x, y});
                
                const point = document.createElement('div');
                point.classList.add('map-point');
                point.style.left = x + 'px';
                point.style.top = y + 'px';
                mapContainer.appendChild(point);
            }
            
            // Calculate midpoint
            const avgX = pointsData.reduce((sum, p) => sum + p.x, 0) / numPoints;
            const avgY = pointsData.reduce((sum, p) => sum + p.y, 0) / numPoints;
            
            // Animate midpoint to new position
            midpoint.style.transition = 'left 1s ease, top 1s ease';
            midpoint.style.left = avgX + 'px';
            midpoint.style.top = avgY + 'px';
        };
        
        // Set up demo buttons
        document.querySelectorAll('.demo-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                document.querySelectorAll('.demo-btn').forEach(b => b.classList.remove('active'));
                this.classList.add('active');
                generatePoints(parseInt(this.dataset.points));
            });
        });
        
        // Initialize with 3 points
        generatePoints(3);
    };
    
    // Initialize demo if we're not on a small screen
    if (window.innerWidth >= 768) {
        createInteractiveDemo();
    }
}); 