<template>

    <!-- HERO BANNER START -->
    <section class="hero-banner">
        <video autoplay="" muted="" loop="" id="myVideo" v-if="homeAnime.video_url">
            <source :src="homeAnime.video_url" type="video/mp4">
        </video>
        <div class="container-fluid">
            <div class="content">
                <div class="txt-block">
                    <div>
                        <p class="white">最新发布</p>
                        <h1 class="mb-12 white">{{ homeAnime.video_name }}</h1>
                        <p class="white mb-48">{{ homeAnime.description }}<br class="d-sm-block d-none">
                            祖安和皮城之间的战争即将展开...</p>
                        <div class="gap-24 d-flex">
                            <a href="" class="cus-btn" @click.prevent="JumpMovies(homeAnime)">
                                <span class="btn-text">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewbox="0 0 20 20"
                                        fill="none">
                                        <path
                                            d="M16.6755 8.8536L4.87079 0.272465C4.43757 -0.0415717 3.86247 -0.0888661 3.38763 0.155174C2.90901 0.397323 2.61011 0.889185 2.61011 1.42078V18.5774C2.61011 19.1146 2.90901 19.6046 3.38763 19.8468C3.59005 19.9489 3.81139 20 4.03462 20C4.32596 20 4.62107 19.9073 4.87079 19.7238L16.6755 11.1502C17.0482 10.8759 17.2658 10.4522 17.2658 10.0019C17.2677 9.5441 17.0444 9.12223 16.6755 8.8536Z"
                                            fill="#FAFAFA"></path>
                                    </svg>
                                    现在播放
                                </span>
                                <span class="btn-text">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewbox="0 0 20 20"
                                        fill="none">
                                        <path
                                            d="M16.6755 8.8536L4.87079 0.272465C4.43757 -0.0415717 3.86247 -0.0888661 3.38763 0.155174C2.90901 0.397323 2.61011 0.889185 2.61011 1.42078V18.5774C2.61011 19.1146 2.90901 19.6046 3.38763 19.8468C3.59005 19.9489 3.81139 20 4.03462 20C4.32596 20 4.62107 19.9073 4.87079 19.7238L16.6755 11.1502C17.0482 10.8759 17.2658 10.4522 17.2658 10.0019C17.2677 9.5441 17.0444 9.12223 16.6755 8.8536Z"
                                            fill="#FAFAFA"></path>
                                    </svg>
                                    现在播放
                                </span>
                            </a>
                            <a href="" class="cus-btn-2" @click.prevent="handleMoreInfo(homeAnime)">
                                <span class="btn-text">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewbox="0 0 20 20"
                                        fill="none">
                                        <g clip-path="url(#clip0_11687_723)">
                                            <path
                                                d="M10 0C4.47254 0 0 4.47293 0 10C0 15.5275 4.47293 20 10 20C15.5275 20 20 15.5271 20 10C20 4.47254 15.527 0 10 0ZM10 18.4375C5.33621 18.4375 1.5625 14.6635 1.5625 10C1.5625 5.33621 5.33652 1.5625 10 1.5625C14.6638 1.5625 18.4375 5.33652 18.4375 10C18.4375 14.6638 14.6634 18.4375 10 18.4375Z"
                                                fill="#FAFAFA"></path>
                                            <path
                                                d="M10 8.37227C9.56852 8.37227 9.21875 8.72203 9.21875 9.15352V14.1845C9.21875 14.616 9.56852 14.9657 10 14.9657C10.4315 14.9657 10.7812 14.6159 10.7812 14.1845V9.15352C10.7812 8.72203 10.4315 8.37227 10 8.37227Z"
                                                fill="#FAFAFA"></path>
                                            <path
                                                d="M10 7.41562C10.5825 7.41562 11.0547 6.94343 11.0547 6.36094C11.0547 5.77845 10.5825 5.30625 10 5.30625C9.41751 5.30625 8.94531 5.77845 8.94531 6.36094C8.94531 6.94343 9.41751 7.41562 10 7.41562Z"
                                                fill="#FAFAFA"></path>
                                        </g>
                                        <defs>
                                            <clippath id="clip0_11687_723">
                                                <rect width="20" height="20" fill="white"></rect>
                                            </clippath>
                                        </defs>
                                    </svg>
                                    更多信息
                                </span>
                                <span class="btn-text">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewbox="0 0 20 20"
                                        fill="none">
                                        <g clip-path="url(#clip0_11687_723)">
                                            <path
                                                d="M10 0C4.47254 0 0 4.47293 0 10C0 15.5275 4.47293 20 10 20C15.5275 20 20 15.5271 20 10C20 4.47254 15.527 0 10 0ZM10 18.4375C5.33621 18.4375 1.5625 14.6635 1.5625 10C1.5625 5.33621 5.33652 1.5625 10 1.5625C14.6638 1.5625 18.4375 5.33652 18.4375 10C18.4375 14.6638 14.6634 18.4375 10 18.4375Z"
                                                fill="#FAFAFA"></path>
                                            <path
                                                d="M10 8.37227C9.56852 8.37227 9.21875 8.72203 9.21875 9.15352V14.1845C9.21875 14.616 9.56852 14.9657 10 14.9657C10.4315 14.9657 10.7812 14.6159 10.7812 14.1845V9.15352C10.7812 8.72203 10.4315 8.37227 10 8.37227Z"
                                                fill="#FAFAFA"></path>
                                            <path
                                                d="M10 7.41562C10.5825 7.41562 11.0547 6.94343 11.0547 6.36094C11.0547 5.77845 10.5825 5.30625 10 5.30625C9.41751 5.30625 8.94531 5.77845 8.94531 6.36094C8.94531 6.94343 9.41751 7.41562 10 7.41562Z"
                                                fill="#FAFAFA"></path>
                                        </g>
                                        <defs>
                                            <clippath id="clip0_11687_724">
                                                <rect width="20" height="20" fill="white"></rect>
                                            </clippath>
                                        </defs>
                                    </svg>
                                    更多信息
                                </span>
                            </a>
                        </div>
                    </div>
                    <button id="togglebutton">
                        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewbox="0 0 48 48" fill="none">
                            <path
                                d="M0.5 24C0.5 11.0213 11.0213 0.5 24 0.5C36.9787 0.5 47.5 11.0213 47.5 24C47.5 36.9787 36.9787 47.5 24 47.5C11.0213 47.5 0.5 36.9787 0.5 24Z"
                                stroke="#FAFAFA"></path>
                            <path
                                d="M15.3522 20.1679H12V27.6321H15.3522L21.6862 32.3296C21.6862 32.3296 22.8479 33.3038 22.8479 32.2981C22.8479 31.2901 22.8479 16.2725 22.8479 15.3559C22.8479 14.566 21.826 15.3374 21.826 15.3374L15.3522 20.1679Z"
                                fill="#FAFAFA"></path>
                            <path
                                d="M30.4765 23.8689L26.3344 18H28.5749L31.6677 22.5246L34.7279 18H36.9464L32.8044 23.8689L37.1759 30H34.9355L31.6677 25.2131L28.3563 30H26.1377L30.4765 23.8689Z"
                                fill="white"></path>
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    </section>
    <!-- HERO BANNER END -->

    <!-- Trending Area Start -->
    <section class="mt-40 mb-48" v-if="japanAnime">
        <div class="container-fluid">
            <div class="heading mb-24">
                <h4 class="white mb-12">日本动漫</h4>
                <router-link to="/animeLibrary?type=japan" class="p white">查看所有 <span><i
                            class="fa-regular fa-chevron-right"></i></span></router-link>
            </div>
            <div class="trendingSlider row">
                <!-- 使用 v-for 渲染列表 -->
                <div class="col-12" v-for="(anime, _) in japanAnime" :key="anime.video_id">
                    <a href="#" class="anime-card gap-24" @click.prevent="JumpMovies(anime)">
                        <img :src="anime.cover_image_url" alt="pic" class="fixed-size-img">
                        <div class="overlay"></div>
                        <div class="text">
                            <h6 class="white">{{ anime.video_name }}</h6>
                            <p class="white">{{ anime.anime_type }}</p>
                        </div>
                    </a>
                </div>
            </div>
        </div>
    </section>
    <!-- Trending Area End -->

    <!-- New Releases Start -->
    <section class="my-48">
        <div class="container-fluid">
            <div class="heading mb-24">
                <h4 class="white mb-12">国产动漫</h4>
                <router-link to="/animeLibrary?type=china" class="p white">查看所有 <span><i
                            class="fa-regular fa-chevron-right"></i></span></router-link>
            </div>
            <div class="releasesSlider row">
                <div class="col-12" v-for="(anime, _) in chinaAnime" :key="anime.video_id">
                    <a href="#" class="anime-card gap-24" @click.prevent="JumpMovies(anime)">
                        <img :src="anime.cover_image_url" alt="pic" class="fixed-size-img">
                        <div class="overlay"></div>
                        <div class="text">
                            <h6 class="white">{{ anime.video_name }}</h6>
                            <p class="white">{{ anime.anime_type }}</p>
                        </div>
                    </a>
                </div>
            </div>
        </div>
    </section>
    <!-- New Releases end -->

    <!-- Top-picks start -->
    <section class="mt-48 mb-40">
        <div class="container-fluid">
            <div class="heading mb-24">
                <h4 class="white mb-12">欧美动漫</h4>
                <router-link to="/animeLibrary?type=western" class="p white">查看所有 <span><i
                            class="fa-regular fa-chevron-right"></i></span></router-link>
            </div>
            <div class="picksSlider row">
                <div class="col-12" v-for="(anime, _) in westernAnime" :key="anime.video_id">
                    <a href="#" class="anime-card gap-24" @click.prevent="JumpMovies(anime)">
                        <img :src="anime.cover_image_url" alt="pic" class="fixed-size-img">
                        <div class="overlay"></div>
                        <div class="text">
                            <h6 class="white">{{ anime.video_name }}</h6>
                            <p class="white">{{ anime.anime_type }}</p>
                        </div>
                    </a>
                </div>
            </div>
        </div>
    </section>
    <!-- Top-picks end -->

    <!-- inner-banner start -->
    <section class="inner-banner">
        <div class="container-fluid">
            <div class="content">
                <ul class="tabs d-flex gap-8 list-unstyled mb-12">
                    <li>热血</li>
                    <li>战斗</li>
                    <li>奇幻</li>
                    <li>国产动漫</li>
                </ul>
                <h1 class="white mb-12">斗破苍穹 第122集</h1>
                <p class="white mb-48">三年之约后，萧炎终于在迦南学院见到了薰儿，此后他广交挚友并成立磐门；为继续提升实力以三上云岚宗为父复仇，<br>
                    他以身犯险深入天焚炼气塔吞噬陨落心炎…影片通过细腻的表现手法，带你进入一个全新的世界。</p>
                <div class="gap-24 d-flex">
                    <a href="" class="cus-btn">
                        <span class="btn-text">
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="21" viewbox="0 0 20 21"
                                fill="none">
                                <g clip-path="url(#clip0_11735_317)">
                                    <path
                                        d="M16.6755 9.3536L4.87079 0.772465C4.43757 0.458428 3.86247 0.411134 3.38763 0.655174C2.90901 0.897323 2.61011 1.38919 2.61011 1.92078V19.0774C2.61011 19.6146 2.90901 20.1046 3.38763 20.3468C3.59005 20.4489 3.81139 20.5 4.03462 20.5C4.32596 20.5 4.62107 20.4073 4.87079 20.2238L16.6755 11.6502C17.0482 11.3759 17.2658 10.9522 17.2658 10.5019C17.2677 10.0441 17.0444 9.62223 16.6755 9.3536Z"
                                        fill="#FAFAFA"></path>
                                </g>
                                <defs>
                                    <clippath id="clip0_11735_317">
                                        <rect width="20" height="20" fill="white" transform="translate(0 0.5)">
                                        </rect>
                                    </clippath>
                                </defs>
                            </svg>
                            观看预告片
                        </span>
                        <span class="btn-text">
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="21" viewbox="0 0 20 21"
                                fill="none">
                                <g clip-path="url(#clip0_11735_317)">
                                    <path
                                        d="M16.6755 9.3536L4.87079 0.772465C4.43757 0.458428 3.86247 0.411134 3.38763 0.655174C2.90901 0.897323 2.61011 1.38919 2.61011 1.92078V19.0774C2.61011 19.6146 2.90901 20.1046 3.38763 20.3468C3.59005 20.4489 3.81139 20.5 4.03462 20.5C4.32596 20.5 4.62107 20.4073 4.87079 20.2238L16.6755 11.6502C17.0482 11.3759 17.2658 10.9522 17.2658 10.5019C17.2677 10.0441 17.0444 9.62223 16.6755 9.3536Z"
                                        fill="#FAFAFA"></path>
                                </g>
                                <defs>
                                    <clippath id="clip0_11735_318">
                                        <rect width="20" height="20" fill="white" transform="translate(0 0.5)">
                                        </rect>
                                    </clippath>
                                </defs>
                            </svg>
                            观看预告片
                        </span>
                    </a>
                    <a href="" class="cus-btn-2">
                        <span class="btn-text">
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="21" viewbox="0 0 20 21"
                                fill="none">
                                <g clip-path="url(#clip0_11735_635)">
                                    <path
                                        d="M16.0547 10.768V8.89844C16.0547 6.17102 14.2418 3.85934 11.7579 3.1043V2.25781C11.7579 1.28855 10.9693 0.5 10 0.5C9.03079 0.5 8.24223 1.28855 8.24223 2.25781V3.1043C5.75825 3.85934 3.94536 6.17098 3.94536 8.89844V10.768C3.94536 13.1638 3.03216 15.4355 1.37399 17.1648C1.21149 17.3342 1.16583 17.5843 1.25794 17.8002C1.35005 18.0161 1.56216 18.1562 1.79692 18.1562H7.1293C7.40149 19.4919 8.58524 20.5 10 20.5C11.4149 20.5 12.5986 19.4919 12.8708 18.1562H18.2032C18.4379 18.1562 18.65 18.0161 18.7421 17.8002C18.8342 17.5843 18.7886 17.3342 18.6261 17.1648C16.9679 15.4355 16.0547 13.1638 16.0547 10.768ZM9.41411 2.25781C9.41411 1.93473 9.67696 1.67188 10 1.67188C10.3231 1.67188 10.586 1.93473 10.586 2.25781V2.87219C10.3931 2.85359 10.1977 2.84375 10 2.84375C9.80239 2.84375 9.60696 2.85359 9.41411 2.87219V2.25781ZM10 19.3281C9.23618 19.3281 8.58489 18.8382 8.34301 18.1562H11.6571C11.4152 18.8382 10.7639 19.3281 10 19.3281ZM3.05981 16.9844C4.39423 15.1956 5.11723 13.0309 5.11723 10.768V8.89844C5.11723 6.20605 7.30766 4.01562 10 4.01562C12.6924 4.01562 14.8829 6.20605 14.8829 8.89844V10.768C14.8829 13.0309 15.6059 15.1956 16.9403 16.9844H3.05981Z"
                                        fill="#FAFAFA"></path>
                                </g>
                                <defs>
                                    <clippath id="clip0_11735_635">
                                        <rect width="20" height="20" fill="white" transform="translate(0 0.5)">
                                        </rect>
                                    </clippath>
                                </defs>
                            </svg>
                            提醒我
                        </span>
                        <span class="btn-text">
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="21" viewbox="0 0 20 21"
                                fill="none">
                                <g clip-path="url(#clip0_11735_636)">
                                    <path
                                        d="M16.0547 10.768V8.89844C16.0547 6.17102 14.2418 3.85934 11.7579 3.1043V2.25781C11.7579 1.28855 10.9693 0.5 10 0.5C9.03079 0.5 8.24223 1.28855 8.24223 2.25781V3.1043C5.75825 3.85934 3.94536 6.17098 3.94536 8.89844V10.768C3.94536 13.1638 3.03216 15.4355 1.37399 17.1648C1.21149 17.3342 1.16583 17.5843 1.25794 17.8002C1.35005 18.0161 1.56216 18.1562 1.79692 18.1562H7.1293C7.40149 19.4919 8.58524 20.5 10 20.5C11.4149 20.5 12.5986 19.4919 12.8708 18.1562H18.2032C18.4379 18.1562 18.65 18.0161 18.7421 17.8002C18.8342 17.5843 18.7886 17.3342 18.6261 17.1648C16.9679 15.4355 16.0547 13.1638 16.0547 10.768ZM9.41411 2.25781C9.41411 1.93473 9.67696 1.67188 10 1.67188C10.3231 1.67188 10.586 1.93473 10.586 2.25781V2.87219C10.3931 2.85359 10.1977 2.84375 10 2.84375C9.80239 2.84375 9.60696 2.85359 9.41411 2.87219V2.25781ZM10 19.3281C9.23618 19.3281 8.58489 18.8382 8.34301 18.1562H11.6571C11.4152 18.8382 10.7639 19.3281 10 19.3281ZM3.05981 16.9844C4.39423 15.1956 5.11723 13.0309 5.11723 10.768V8.89844C5.11723 6.20605 7.30766 4.01562 10 4.01562C12.6924 4.01562 14.8829 6.20605 14.8829 8.89844V10.768C14.8829 13.0309 15.6059 15.1956 16.9403 16.9844H3.05981Z"
                                        fill="#FAFAFA"></path>
                                </g>
                                <defs>
                                    <clippath id="clip0_11735_637">
                                        <rect width="20" height="20" fill="white" transform="translate(0 0.5)">
                                        </rect>
                                    </clippath>
                                </defs>
                            </svg>
                            提醒我
                        </span>
                    </a>
                </div>
            </div>
        </div>
    </section>
    <!-- inner-banner end -->

    <!--accordion start -->
    <section class=" mt-48 mb-24">
        <div class="container-fluid">
            <div class="heading mb-24">
                <h4 class="white mb-12">为你推荐</h4>
                <router-link to="/animeLibrary?type=top" class="p white">查看所有 <span><i
                            class="fa-regular fa-chevron-right"></i></span></router-link>
            </div>
            <div class="anime-grid">
                <a href="#" class="anime-card gap-24" v-for="(anime, _) in topAnime" :key="anime.video_id"
                    @click.prevent="JumpMovies(anime)">
                    <img :src="anime.cover_image_url" alt="pic" class="fixed-size-img">
                    <div class="overlay"></div>
                    <div class="text">
                        <h6 class="white">{{ anime.video_name }}</h6>
                        <!-- <p class="white">Mystery, Supernatural</p> -->
                    </div>
                </a>
            </div>
            <div class="accordion" id="accordionPanelsStayOpenExample">
                <div class="accordion-item">
                    <h2 class="accordion-header" id="panelsStayOpen-headingOne">
                        <button class="accordion-button" type="button" data-bs-toggle="collapse"
                            data-bs-target="#panelsStayOpen-collapseOne" aria-expanded="true"
                            aria-controls="panelsStayOpen-collapseOne">
                            查看更多 <span><i class="fa-light fa-chevron-down"></i></span>
                        </button>
                    </h2>
                    <div id="panelsStayOpen-collapseOne" class="accordion-collapse collapse"
                        aria-labelledby="panelsStayOpen-headingOne">
                        <div class="accordion-body">
                            <div class="anime-grid">
                                <a href="animeDetail" class="anime-card gap-24">
                                    <img src="@/static/picture/anime-01.jpg" alt="pic">
                                    <div class="overlay"></div>
                                    <div class="text">
                                        <h6 class="white">Death Note</h6>
                                        <p class="white">Mystery, Supernatural</p>
                                    </div>
                                </a>
                                <a href="animeDetail" class="anime-card gap-24">
                                    <img src="@/static/picture/anime-02.jpg" alt="pic">
                                    <div class="overlay"></div>
                                    <div class="text">
                                        <h6 class="white">Fullmetal Alchemist</h6>
                                        <p class="white">Action, Adventure, Steampunk</p>
                                    </div>
                                </a>
                                <a href="animeDetail" class="anime-card gap-24">
                                    <img src="@/static/picture/anime-03.jpg" alt="pic">
                                    <div class="overlay"></div>
                                    <div class="text">
                                        <h6 class="white">Sword Art Online</h6>
                                        <p class="white">Action, Adventure, Fantasy</p>
                                    </div>
                                </a>
                                <a href="animeDetail" class="anime-card gap-24">
                                    <img src="@/static/picture/anime-04.jpg" alt="pic">
                                    <div class="overlay"></div>
                                    <div class="text">
                                        <h6 class="white">Tokyo Ghoul</h6>
                                        <p class="white">Action, Horror, Supernatural</p>
                                    </div>
                                </a>
                                <a href="animeDetail" class="anime-card gap-24">
                                    <img src="@/static/picture/anime-05.jpg" alt="pic">
                                    <div class="overlay"></div>
                                    <div class="text">
                                        <h6 class="white">Komitsu Uhang</h6>
                                        <p class="white">Anime, Vervex</p>
                                    </div>
                                </a>
                                <a href="animeDetail" class="anime-card d-xxl-none gap-24">
                                    <img src="@/static/picture/anime-06.jpg" alt="pic">
                                    <div class="overlay"></div>
                                    <div class="text">
                                        <h6 class="white">Komitsu Uhang</h6>
                                        <p class="white">Anime, Vervex</p>
                                    </div>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    </section>
    <!--accordion start  -->

    <!-- collection start -->
    <section class="collection my-24">
        <div class="container-fluid">
            <h4 class="mb-24 white">动漫类型</h4>
            <div class="row row-gap-3 justify-content-center">
                <div class="col-xl-3 col-md-4 col-sm-6" v-for="(anime, index) in animeGenres" :key="index">
                    <router-link to="animeLibrary?type=genre" class="cards">
                        <h5 class="white">{{ anime }}</h5>
                        <img :src="'/src/static/picture/icon' + (index + 1) + '.png'" alt="icon">
                    </router-link>
                </div>
            </div>
        </div>
    </section>
    <!-- collection end -->
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            homeAnime: [],
            japanAnime: [],
            chinaAnime: [],
            westernAnime: [],
            topAnime: [],
            animeGenres: [],
            notifications: [
                { id: 1, title: "新番更新！", message: "《鬼灭之刃》最新一集已上线，快来观看！", time: "2分钟前" },
                { id: 2, title: "漫画更新", message: "《海贼王》新章节已发布，立即查看！", time: "15分钟前" },
                { id: 3, title: "活动提醒", message: "下周是漫展最后一天，别错过哦！", time: "1小时前" },
            ],
        };
    },
    methods: {
        JumpMovies(result) {
            this.$router.push(`/moviesDetail?videoId=${result.video_id}`)
        },

        async fetchHomeAnime() {
            try {
                const response = await axios.get('/api/getHomeAnime')
                if (response.data.code == 200) {
                    this.homeAnime = response.data.home_anime
                    this.japanAnime = response.data.japan_anime
                    this.chinaAnime = response.data.china_anime
                    this.westernAnime = response.data.western_anime
                    this.topAnime = response.data.top_anime
                    this.animeGenres = response.data.anime_genres
                } else {
                    console.error("fetchHomeAnime:", response.data.message)
                }
            }
            catch (error) {
                console.error("fetchHomeAnime:", error.response.data)
            }
        },

        handleMoreInfo(info) {
            this.$router.push(`/animeDetail?params=${encodeURIComponent(info.video_name.trim())}`)
        },

        initEvent() {
            if ($("#myVideo").length) {
                const video = $("#myVideo")[0];
                const toggleButton = $("#togglebutton")[0];
                let isMuted = true;
                // Initially, video is muted
                video.play();

                toggleButton.addEventListener("click", function () {
                    isMuted = !isMuted;
                    video.muted = isMuted;
                    if (isMuted) {
                        $('#togglebutton svg').html('<path d="M0.5 24C0.5 11.0213 11.0213 0.5 24 0.5C36.9787 0.5 47.5 11.0213 47.5 24C47.5 36.9787 36.9787 47.5 24 47.5C11.0213 47.5 0.5 36.9787 0.5 24Z" stroke="#FAFAFA"/> <path d="M15.3522 20.1679H12V27.6321H15.3522L21.6862 32.3296C21.6862 32.3296 22.8479 33.3038 22.8479 32.2981C22.8479 31.2901 22.8479 16.2725 22.8479 15.3559C22.8479 14.566 21.826 15.3374 21.826 15.3374L15.3522 20.1679Z" fill="#FAFAFA"/> <path d="M30.4765 23.8689L26.3344 18H28.5749L31.6677 22.5246L34.7279 18H36.9464L32.8044 23.8689L37.1759 30H34.9355L31.6677 25.2131L28.3563 30H26.1377L30.4765 23.8689Z" fill="white"/>');
                    } else {
                        $('#togglebutton svg').html('<path d="M0.5 24C0.5 11.0213 11.0213 0.5 24 0.5C36.9787 0.5 47.5 11.0213 47.5 24C47.5 36.9787 36.9787 47.5 24 47.5C11.0213 47.5 0.5 36.9787 0.5 24Z" stroke="#FAFAFA" /> <path d="M15.2145 20.6442H11.8623V28.1085H15.2145L21.5485 32.806C21.5485 32.806 22.7102 33.7802 22.7102 32.7745C22.7102 31.7664 22.7102 16.7489 22.7102 15.8323C22.7102 15.0424 21.6883 15.8138 21.6883 15.8138L15.2145 20.6442Z" fill="#FAFAFA" /> <path d="M26.319 19.1875C25.984 18.854 25.4454 18.854 25.1143 19.1875C24.7793 19.5209 24.7793 20.0618 25.1143 20.3922C26.2099 21.4917 26.7531 22.9169 26.7546 24.3544C26.7531 25.795 26.2099 27.2295 25.1143 28.3236C24.7793 28.657 24.7793 29.1964 25.1143 29.5321C25.2795 29.6996 25.4969 29.7811 25.7174 29.7811C25.9333 29.7811 26.1523 29.6996 26.319 29.5321C27.7442 28.1038 28.4603 26.2237 28.4588 24.3544C28.4611 22.4836 27.745 20.6135 26.319 19.1875Z" fill="#FAFAFA" /> <path d="M28.6301 16.5022C28.2936 16.1672 27.755 16.1672 27.4223 16.5022C27.0919 16.8364 27.0919 17.3758 27.4223 17.71C29.2647 19.5509 30.1821 21.9511 30.1852 24.3614C30.1821 26.7846 29.2686 29.1956 27.4223 31.0473C27.0904 31.3807 27.0919 31.9178 27.4223 32.2535C27.5906 32.4164 27.808 32.5009 28.027 32.5009C28.2444 32.5009 28.4634 32.4164 28.6301 32.2535C30.806 30.0715 31.8893 27.2118 31.8893 24.3614C31.8893 21.5201 30.7975 18.672 28.6301 16.5022Z" fill="#FAFAFA" /> <path d="M31.4606 13.7516C31.1287 13.4166 30.5885 13.4166 30.2543 13.7516C29.9239 14.0851 29.9239 14.626 30.2543 14.9563C32.8566 17.5586 34.1558 20.9561 34.1558 24.3644C34.1558 27.7827 32.8604 31.1902 30.2543 33.7978C29.9209 34.1305 29.9239 34.673 30.2543 35.0041C30.4226 35.1693 30.6415 35.2515 30.8605 35.2515C31.0764 35.2515 31.2938 35.1693 31.4606 35.0041C34.3971 32.0653 35.863 28.2083 35.8623 24.3636C35.863 20.5259 34.3894 16.6804 31.4606 13.7516Z" fill="#FAFAFA" />');
                    }
                });
            }
        }
    },
    async mounted() {
        await this.fetchHomeAnime()
        app()
        this.initEvent()
    },
};

import app from '@/static/js/app';
</script>

<style scoped>
.fixed-size-img {
    width: 400px;
    /* 设置你想要的宽度 */
    height: 250px;
    /* 设置你想要的高度 */
    object-fit: cover;
    /* 保持图片比例，裁剪多余部分 */
}

.input-box {
    position: relative;
}

.search-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: #1A1A1A;
    border: 1px solid #333;
    border-radius: 4px;
    margin-top: 4px;
    max-height: 300px;
    overflow-y: auto;
    z-index: 1000;
}

.search-dropdown ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.search-result-item {
    display: flex;
    align-items: center;
    padding: 10px;
    cursor: pointer;
    transition: background-color 0.2s;
}

.search-result-item:hover {
    background-color: #2A2A2A;
}

.result-thumbnail {
    width: 40px;
    height: 40px;
    object-fit: cover;
    border-radius: 4px;
    margin-right: 10px;
}

.result-info {
    flex: 1;
}

.result-title {
    color: #FFFFFF;
    font-size: 14px;
    margin-bottom: 4px;
}

.result-type {
    color: #999;
    font-size: 12px;
}

/* 自定义滚动条样式 */
.search-dropdown::-webkit-scrollbar {
    width: 6px;
}

.search-dropdown::-webkit-scrollbar-track {
    background: #1A1A1A;
}

.search-dropdown::-webkit-scrollbar-thumb {
    background: #444;
    border-radius: 3px;
}

.search-dropdown::-webkit-scrollbar-thumb:hover {
    background: #555;
}
</style>