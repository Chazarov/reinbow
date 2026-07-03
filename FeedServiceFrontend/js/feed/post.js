window.prevImage = (postId) => {
    const imagesBlock = document.getElementById("images_" + postId);
    const currentImage = imagesBlock.querySelector(".showedImage");
    const prevImage = currentImage.previousElementSibling;
    if(prevImage) {
        currentImage.classList.remove('showedImage');
        prevImage.classList.add('showedImage');
    }
}

window.nextImage = (postId) => {
    const imagesBlock = document.getElementById("images_" + postId);
    const currentImage = imagesBlock.querySelector(".showedImage");
    const nextImage = currentImage.nextElementSibling;
    if(nextImage) {
        currentImage.classList.remove('showedImage');
        nextImage.classList.add('showedImage');
    }
}


export function createPost(post) {
    const imageUrl = process.env.IS_URL;   
    
    let sliderImages = "<div>"

    // todo в src вставить ${process.env.IS_URL}/api/images/${imageId}"
    //const test = 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSPshUBFBR7xshIaXFc_ir-eagtAueBv7aX5Nvxdny6qg&s'

    for(const imageId of post.images || []) {
        sliderImages += `<img id="image_${imageId}" src="${process.env.IS_URL}/api/images/${imageId}"" class="postImage">`
    }
    sliderImages += "</div>"

    const tools = `
        <div class="imageTools">
            <span onClick="prevImage(${post.id})" class="imageBtn imagePrev">&larr;</span>
            <span onClick="nextImage(${post.id})" class="imageBtn imageNext">&rarr;</span>
        </div>`


    return `
        <div class="post-card" id="post_${post.id}">
            <div class="imagesSlider">
                <div class="images" id="images_${post.id}">
                    ${sliderImages}
                    ${post.images && post.images.length > 1 ? tools : ''}
                </div>
            </div>
            <div class="post-body">
                <h3 class="post-title">${post.title || 'Без названия'}</h3>
                <p class="post-description">${post.description || ''}</p>
                <div class="post-meta">
                    <span class="post-author" data-user-id="${post.user_id}">${post.username || 'Аноним'}</span>
                    <span class="post-date">${post.created_at}</span>
                </div>
            </div>
        </div>
    `;
}

export function createPostsHandlers() {
    document.querySelectorAll('.postImage:first-child').forEach(el => {
        el.classList.add("showedImage")
    });

    document.querySelectorAll('.post-author').forEach(el => {
        el.addEventListener('click', (e) => {
            const userId = e.target.dataset.userId;
            if (userId) {
                loadUserFeed(parseInt(userId));
            }
        });
    });
}