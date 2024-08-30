const menuItems = document.querySelectorAll('.menu-item');
const votingContent = document.getElementById('voting-content');
const breedsContent = document.getElementById('breeds-content');
const favsContent = document.getElementById('favs-content');
const gridViewButton = document.getElementById('grid-view');
const listViewButton = document.getElementById('list-view');
const favoritesGrid = document.getElementById('favoritesGrid');

let votes = 0;
let favorites = 0;
let currentImageId = '';
let breeds = [];

function showContent(contentId) {
    [votingContent, breedsContent, favsContent].forEach(content => {
        content.style.display = content.id === contentId ? 'block' : 'none';
    });
}

function showBreeds() {
    document.getElementById('breedSelect').style.display = 'block';
}

function showFavs() {
    document.getElementById('favoritesSection').style.display = 'block';
    fetchFavorites();
    //console.log('Favs clicked');
}

menuItems.forEach(item => {
    item.addEventListener('click', function() {
        menuItems.forEach(i => i.classList.remove('active'));
        this.classList.add('active');
        
        if (this.id === 'voting-menu') {
            showContent('voting-content');
        } else if (this.id === 'breeds-menu') {
            showContent('breeds-content');
            showBreeds(); // Optionally call this if needed
        } else if (this.id === 'favs-menu') {
            showContent('favs-content');
            showFavs(); // Ensure this is called
        }
    });
});

// Set Voting as default active view
document.getElementById('voting-menu').click();

function fetchCatData() {
    fetch('/api/cat')
        .then(response => response.json())
        .then(data => {
            document.getElementById('catImage').src = data.url;
            document.getElementById('catBreed').textContent = data.breeds && data.breeds.length > 0 ? data.breeds[0].name : 'Unknown';
            currentImageId = data.id;
        })
        .catch(error => console.error('Error fetching cat data:', error));
}

function handleFavorite() {
    fetch('/api/favorites', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ image_id: currentImageId }),
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(errorData => {
                throw new Error(errorData.error || 'Unknown error');
            });
        }
        return response.json();
    })
    .then(data => {
        if (data.id) {
            //alert('Added to favorites! Favorite ID: ' + data.id);
            fetchCatData();  // Fetch a new cat image
        } else {
            //alert('Unexpected response from server');
        }
    })
    .catch(error => {
        console.error('Error adding favorite:', error);
        //alert('Error adding favorite: ' + error.message);
    });

}

function fetchFavorites() {
    fetch('/api/favorites')
        .then(response => response.json())
        .then(data => {
            const favoritesGrid = document.getElementById('favoritesGrid');
            favoritesGrid.innerHTML = '';
            data.forEach(favorite => {
                const img = document.createElement('img');
                img.src = favorite.image.url;
                img.alt = 'Favorite cat';
                favoritesGrid.appendChild(img);
            });
        })
        .catch(error => console.error('Error fetching favorites:', error));
}

function submitVote(value) {
    if (!currentImageId) {
        console.error('No image to vote on');
        return;
    }

    const voteData = {
        image_id: currentImageId,
        sub_id: "user-" + Date.now(), // Generate a unique sub_id
        value: value
    };

    fetch('/api/vote', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(voteData)
    })
    .then(response => response.json())
    .then(data => {
        if (data.message === "SUCCESS") {
            fetchCatData();
            //alert(`Vote submitted successfully! ${value > 0 ? 'Upvote' : 'Downvote'}`);
        } else {
            //alert('Failed to submit vote. Please try again.');
        }
    })
    .catch(error => {
        console.error('Error submitting vote:', error);
        alert('An error occurred while submitting your vote.');
    });
}

function fetchBreeds() {
    fetch('/api/breeds')
        .then(response => response.json())
        .then(data => {
            breeds = data;
            populateBreedSelect();
            if (breeds.length > 0) {
                fetchBreedInfo(breeds[0].id);
            }
        })
        .catch(error => console.error('Error fetching breeds:', error));
}

function populateBreedSelect() {
    const select = document.getElementById('breedSelect');
    select.innerHTML = '';
    breeds.forEach(breed => {
        const option = document.createElement('option');
        option.value = breed.id;
        option.textContent = breed.name;
        select.appendChild(option);
    });
    select.addEventListener('change', (event) => fetchBreedInfo(event.target.value));
}

function fetchBreedInfo(breedId) {
    fetch(`/api/breed/${breedId}`)
        .then(response => response.json())
        .then(data => {
            currentBreed = data.breed_info;
            updateBreedDisplay(data);
        })
        .catch(error => console.error('Error fetching breed info:', error));
}

function updateBreedDisplay(data) {
    document.getElementById('breedName').textContent = currentBreed.name;
    document.getElementById('breedId').textContent = currentBreed.id;
    document.getElementById('breedDescription').textContent = currentBreed.description;
    document.getElementById('breedOrigin').textContent = currentBreed.origin;
    document.getElementById('wikiLink').href = currentBreed.wikipedia_url;

    const sliderContainer = document.getElementById('sliderImagesContainer');
    const navDots = document.getElementById('sliderNavDots');
    sliderContainer.innerHTML = '';  // Clear existing images
    navDots.innerHTML = '';          // Clear existing navigation dots

    data.images.forEach((image, index) => {
        const img = document.createElement('img');
        img.src = image.url;
        img.alt = `Breed Cat image ${index + 1}`;
        img.classList.add('slider-image');
        if (index === 0) img.classList.add('active'); // Make the first image active initially
        sliderContainer.appendChild(img);

        const dot = document.createElement('div');
        dot.classList.add('dot');
        if (index === 0) dot.classList.add('active');
        dot.addEventListener('click', () => goToSlide(index));
        navDots.appendChild(dot);
    });

    let currentIndex = 0;
    const images = document.querySelectorAll('.slider-image');
    const dots = document.querySelectorAll('.dot');

    function goToSlide(index) {
        images[currentIndex].classList.remove('active');
        dots[currentIndex].classList.remove('active');

        currentIndex = index;

        images[currentIndex].classList.add('active');
        dots[currentIndex].classList.add('active');
    }

    // Auto-advance slides every 5 seconds
    setInterval(() => {
        goToSlide((currentIndex + 1) % images.length);
    }, 5000);
}

gridViewButton.addEventListener('click', () => {
    favoritesGrid.classList.remove('single-image-scroll');
    favoritesGrid.classList.add('favorites-grid');
});

listViewButton.addEventListener('click', () => {
    favoritesGrid.classList.remove('favorites-grid');
    favoritesGrid.classList.add('single-image-scroll');
});

// Initial load
fetchCatData();
fetchBreeds();