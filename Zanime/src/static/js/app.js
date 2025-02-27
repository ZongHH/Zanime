// =================== ****************** ================== //

import axios from "axios";
import { ElMessage } from 'element-plus'

// Template Name: Anime tube
// Description:  Anime tube
// Version: 1.0.0

// =================== ****************** ================== //

var MyScroll = "";
var Init = {
  i: function (e) {
    Init.s();
    Init.methods();
  },
  s: function (e) {
    (this._window = $(window)),
      (this._document = $(document)),
      (this._body = $("body")),
      (this._html = $("html"));
  },
  methods: function (e) {
    Init.preloader();
    Init.backToTop();
    Init.header();
    Init.categoryToggle();
    Init.filterSearch();
    Init.passwordIcon();
    Init.slickSlider();
    Init.searchFunction();
    Init.videoFunction();
    Init.videoPlay();
    Init.VideoPlayer();
    Init.formValidation();
    Init.niceSelect();
  },

  showError: function (message) {
    ElMessage.error("服务器繁忙，请刷新重试")
  },

  initVideoPlayerFromBackend: async function () {
    // 获取当前页面 URL 中的 videoId 参数
    const urlParams = new URLSearchParams(window.location.search);
    const videoId = urlParams.get('videoId');
    const episode = urlParams.get('episode');

    try {
      const response = await axios.get(`/api/video-resource?videoId=${videoId}&episode=${episode}`)
      if (response.data.code == 200) {
        Init.initVideoPlayer("#video", response.data.VideoFiles, response.data.PosterPath, true);
      } else {
        Init.showError(response.data.message)
      }

    } catch (error) {
      Init.showError(error)
    }
  },

  initVideoPlayer: function (videoElementId, videoFiles, posterPath, forwardEnabled) {
    if ($(videoElementId).length) {
      $(videoElementId).aksVideoPlayer({
        file: videoFiles,
        poster: posterPath,
        forward: forwardEnabled
      });
    }
  },

  // search
  searchFunction: function () {
    let inputBox = document.querySelector(".input-box"),
      searchIcon = document.querySelector(".search"),
      closeIcon = document.querySelector(".close-icon");
    searchIcon.addEventListener("click", () => {
      inputBox.classList.add("open");
    });
    closeIcon.addEventListener("click", () => {
      inputBox.classList.remove("open");
    });
  },

  // Preloader
  preloader: function () {
    setTimeout(function () { $('#preloader').hide('slow') }, 2000);
  },
  //backtotop
  backToTop: function () {
    var scrollToTopBtn = document.querySelector(".scrollToTopBtn");
    var rootElement = document.documentElement;
    function handleScroll() {
      var scrollTotal = rootElement.scrollHeight - rootElement.clientHeight;
      if (rootElement.scrollTop / scrollTotal > 0.05) {
        scrollToTopBtn.classList.add("showBtn");
      } else {
        scrollToTopBtn.classList.remove("showBtn");
      }
    }
    function scrollToTop() {
      rootElement.scrollTo({ top: 0, behavior: "smooth" });
    }
    scrollToTopBtn.addEventListener("click", scrollToTop);
    document.addEventListener("scroll", handleScroll);
  },
  // Header
  header: function () {
    function dynamicCurrentMenuClass(selector) {
      let FileName = window.location.href.split("/").reverse()[0];

      selector.find("li").each(function () {
        let anchor = $(this).find("a");
        if ($(anchor).attr("href") == FileName) {
          $(this).addClass("current");
        }
      });
      selector.children("li").each(function () {
        if ($(this).find(".current").length) {
          $(this).addClass("current");
        }
      });
      if ("" == FileName) {
        selector.find("li").eq(0).addClass("current");
      }
    }

    if ($(".main-menu__list").length) {
      let mainNavUL = $(".main-menu__list");
      dynamicCurrentMenuClass(mainNavUL);
    }

    if ($(".main-menu__nav").length && $(".mobile-nav__container").length) {
      let navContent = document.querySelector(".main-menu__nav").innerHTML;
      let mobileNavContainer = document.querySelector(
        ".mobile-nav__container"
      );
      mobileNavContainer.innerHTML = navContent;
    }
    if ($(".sticky-header__content").length) {
      let navContent = document.querySelector(".main-menu").innerHTML;
      let mobileNavContainer = document.querySelector(
        ".sticky-header__content"
      );
      mobileNavContainer.innerHTML = navContent;
    }

    if ($(".mobile-nav__container .main-menu__list").length) {
      let dropdownAnchor = $(
        ".mobile-nav__container .main-menu__list .dropdown > a"
      );
      dropdownAnchor.each(function () {
        let self = $(this);
        let toggleBtn = document.createElement("BUTTON");
        toggleBtn.setAttribute("aria-label", "dropdown toggler");
        toggleBtn.innerHTML = "<i class='fa fa-angle-down'></i>";
        self.append(function () {
          return toggleBtn;
        });
        self.find("button").on("click", function (e) {
          e.preventDefault();
          let self = $(this);
          self.toggleClass("expanded");
          self.parent().toggleClass("expanded");
          self.parent().parent().children("ul").slideToggle();
        });
      });
    }

    if ($(".mobile-nav__toggler").length) {
      $(".mobile-nav__toggler").on("click", function (e) {
        e.preventDefault();
        $(".mobile-nav__wrapper").toggleClass("expanded");
        $("body").toggleClass("locked");
      });
    }

    $(window).on("scroll", function () {
      if ($(".stricked-menu").length) {
        var headerScrollPos = 130;
        var stricky = $(".stricked-menu");
        if ($(window).scrollTop() > headerScrollPos) {
          stricky.addClass("stricky-fixed");
        } else if ($(this).scrollTop() <= headerScrollPos) {
          stricky.removeClass("stricky-fixed");
        }
      }
    });
  },
  niceSelect: function () {
    if ($(".has-nice-select").length) {
      $(".has-nice-select").niceSelect();
    }
  },

  // Filter Toggle Button
  categoryToggle: function () {
    if ($(".category-block").length) {
      $(".category-block .title").on("click", function (e) {
        var count = $(this).data("count");
        if (
          $(".category-block.box-" + count + " .content-block").is(":visible")
        ) {
          $(".category-block.box-" + count + " i").removeClass(
            "fa-horizontal-rule"
          );
          $(".category-block.box-" + count + " i").addClass("fa-plus");
          $(".category-block.box-" + count + " .content-block").hide("slow");
        } else {
          $(".category-block.box-" + count + " i").removeClass("fa-plus");
          $(".category-block.box-" + count + " i").addClass(
            "fa-horizontal-rule"
          );
          $(".category-block.box-" + count + " .content-block").show("slow");
        }
      });
    }
    if ($(".customer-container").length) {
      $(".signin-button").click(function () {
        $(".sign-form").slideToggle();
      });
    }

    if ($("#shipAddress").length) {
      $(".billing-address").hide();
      $("#shipAddress").change(function () {
        if ($(this).is(":unchecked")) {
          $(".billing-address").hide("slow");
        } else {
          $(".billing-address").show("slow");
        }
      });
    }
  },

  filterSearch: function () {
    if ($("#searchInput").length) {
      $("#searchInput").on("keyup", function () {
        var value = $(this).val().toLowerCase();
        $(".blogs-block").filter(function () {
          var hasMatch =
            $(this).find(".blog-title").text().toLowerCase().indexOf(value) >
            -1;
          $(this).toggle(hasMatch);
        });
      });
    }
  },

  // Password Icon
  togglepassword2: function () {
    function togglePassword() {
      const passwordInput = document.getElementById("password");
      const currentType = passwordInput.getAttribute("type");
      passwordInput.setAttribute(
        "type",
        currentType === "password" ? "text" : "password"
      );
    }
  },

  //  slider
  slickSlider: function () {
    if ($(".trendingSlider ").length) {
      $(".trendingSlider ").slick({
        infinite: true,
        slidesToShow: 5,
        slidesToScroll: 1,
        arrows: false,
        centerPadding: "0",
        autoplay: false,
        autoplaySpeed: 3000,
        responsive: [
          {
            breakpoint: 1480,
            settings: {
              slidesToShow: 3,
            },
          },
          {
            breakpoint: 1199,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 768,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 390,
            settings: {
              slidesToShow: 1,
            },
          },
        ],
      });
    }
    if ($(".releasesSlider ").length) {
      $(".releasesSlider ").slick({
        infinite: true,
        slidesToShow: 5,
        slidesToScroll: 1,
        arrows: false,
        centerPadding: "0",
        autoplay: false,
        autoplaySpeed: 3000,
        responsive: [
          {
            breakpoint: 1480,
            settings: {
              slidesToShow: 3,
            },
          },
          {
            breakpoint: 1199,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 768,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 390,
            settings: {
              slidesToShow: 1,
            },
          },
        ],
      });
    }
    if ($(".picksSlider").length) {
      $(".picksSlider").slick({
        infinite: true,
        slidesToShow: 5,
        slidesToScroll: 1,
        arrows: false,
        centerPadding: "0",
        autoplay: false,
        autoplaySpeed: 3000,
        responsive: [
          {
            breakpoint: 1480,
            settings: {
              slidesToShow: 3,
            },
          },
          {
            breakpoint: 1199,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 768,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 390,
            settings: {
              slidesToShow: 1,
            },
          },
        ],
      });
    }
    if ($(".blogsSlider").length) {
      $(".blogsSlider").slick({
        infinite: true,
        slidesToShow: 3,
        slidesToScroll: 1,
        arrows: false,
        centerPadding: "0",
        autoplay: false,
        autoplaySpeed: 3000,
        responsive: [
          {
            breakpoint: 1480,
            settings: {
              slidesToShow: 3,
            },
          },
          {
            breakpoint: 1199,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 768,
            settings: {
              slidesToShow: 2,
            },
          },
          {
            breakpoint: 577,
            settings: {
              slidesToShow: 1,
            },
          },
          {
            breakpoint: 390,
            settings: {
              slidesToShow: 1,
            },
          },
        ],
      });
    }
    $(".btn-prev").click(function () {
      var $this = $(this).attr("data-slide");
      $("." + $this).slick("slickPrev");
    });
    $(".btn-next").click(function () {
      var $this = $(this).attr("data-slide");
      $("." + $this).slick("slickNext");
    });
  },

  //  video function
  videoFunction: function () {
    if ($("#myVideo").length) {
      document.addEventListener("DOMContentLoaded", function () {
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
      });
    }
  },
  VideoPlayer: function () {
    if (document.querySelector("#video") !== null) {
      Init.initVideoPlayerFromBackend();
    }
  },
  videoPlay: function () {
    $('#videoModal').on('hidden.bs.modal', function () {
      $('#videoModal video').get(0).pause();
    });
    $("#closeVideoModalButton").click(function () {
      $("#videoModal").modal("hide");
    });
  },
  // Countdown Timer
  countdownInit: function (countdownSelector, countdownTime, countdown) {
    var eventCounter = $(countdownSelector);
    if (eventCounter.length) {
      eventCounter.countdown(countdownTime, function (e) {
        $(this).html(
          e.strftime(
            "<li><h2>%D</h2><h6>Days</h6></li>\
            <li><h2>%H</h2><h6>Hrs</h6></li>\
            <li><h2>%M</h2><h6>Mins</h6></li>\
            <li><h2>%S</h2><h6>Secs</h6></li>"
          )
        );
      });
    }
  },

  // Form Validation
  formValidation: function () {
    if ($(".contact-form").length) {
      $(".contact-form").validate();
    }

    if ($(".login-form").length) {
      $(".login-form").validate();
    }
  },
  passwordIcon: function () {
    $("#eye , #eye-icon").click(function () {
      if ($(this).hasClass("fa-eye-slash")) {
        $(this).removeClass("fa-eye-slash");
        $(this).addClass("fa-eye");
        $(".password-input").attr("type", "text");
      } else {
        $(this).removeClass("fa-eye");
        $(this).addClass("fa-eye-slash");
        $(".password-input").attr("type", "password");
      }
    });
  },

  // Contact Form
  contactForm: function () {
    $(".contact-form").on("submit", function (e) {
      e.preventDefault();
      if ($(".contact-form").valid()) {
        var _self = $(this);
        _self
          .closest("div")
          .find('button[type="submit"]')
          .attr("disabled", "disabled");
        var data = $(this).serialize();
        $.ajax({
          url: "./assets/mail/contact.php",
          type: "post",
          dataType: "json",
          data: data,
          success: function (data) {
            $(".contact-form").trigger("reset");
            _self.find('button[type="submit"]').removeAttr("disabled");
            if (data.success) {
              document.getElementById("message").innerHTML =
                "<h4 class='color-primary mt-16 mb-16'>Email Sent Successfully</h4>";
            } else {
              document.getElementById("message").innerHTML =
                "<h4 class='color-primary mt-16 mb-16'>There is an error</h4>";
            }
            $("#messages").show("slow");
            $("#messages").slideDown("slow");
            setTimeout(function () {
              $("#messages").slideUp("hide");
              $("#messages").hide("slow");
            }, 4000);
          },
        });
      } else {
        return false;
      }
    });
  },
  bookingForm: function () {
    $(".bookingProcced").on("click", function () {
      $(".stepone").hide("slow");
      $(".steptwo").show("slow");
    });
  },
};

export default async function app() {
  "use strict";
  // ==========================================================
  // Detect mobile device and add class "is-mobile" to </body>
  // ==========================================================

  // Detect mobile device (Do not remove!!!)
  var isMobile =
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Nokia|Opera Mini/i.test(
      navigator.userAgent
    )
      ? true
      : false;
  var Scrollbar = window.Scrollbar;
  // Add class "is-mobile" to </body>

  Init.i();
}

// 修改导出的 VideoPlayer 函数
export async function VideoPlayer() {
  // 获取当前页面 URL 中的 videoId 参数
  const urlParams = new URLSearchParams(window.location.search);
  const videoId = urlParams.get('videoId');
  const episode = urlParams.get('episode');

  try {
    const response = await fetch(`/api/video-resource?videoId=${videoId}&episode=${episode}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    console.log(data.VideoFiles);

    // 直接调用初始化视频播放器
    if (document.querySelector("#video")) {
      $("#video").aksVideoPlayer({
        file: data.VideoFiles,
        poster: data.PosterPath,
        forward: true
      });
    }
  } catch (error) {
    Init.showError(error);
  }
}

export async function clearVideoPlayer() {
  // 找到并移除所有包含 aks-video-player 类的元素
  const videoPlayers = document.getElementsByClassName('aks-video-player');
  while (videoPlayers.length > 0) {
    videoPlayers[0].remove();
  }
}

export async function clearVideoPlayerOLD() {
  // 找到并移除所有包含 aks-video-player 类的元素，只保留一个
  const videoPlayers = document.getElementsByClassName('aks-video-player');
  while (videoPlayers.length > 1) {
    videoPlayers[videoPlayers.length - 1].remove();
  }
}