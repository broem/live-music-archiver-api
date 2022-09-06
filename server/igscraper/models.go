package igscraper

import "encoding/json"

type FeedMedia struct {
	err error

	uid       int64
	endpoint  string
	timestamp string

	Items               []Item      `json:"items"`
	NumResults          int         `json:"num_results"`
	MoreAvailable       bool        `json:"more_available"`
	AutoLoadMoreEnabled bool        `json:"auto_load_more_enabled"`
	Status              string      `json:"status"`
	NextID              interface{} `json:"next_max_id"`
}

type Item struct {
	Comments *Comments `json:"-"`

	TakenAt          int64   `json:"taken_at"`
	Pk               int64   `json:"pk"`
	ID               string  `json:"id"`
	CommentsDisabled bool    `json:"comments_disabled"`
	DeviceTimestamp  int64   `json:"device_timestamp"`
	MediaType        int     `json:"media_type"`
	Code             string  `json:"code"`
	ClientCacheKey   string  `json:"client_cache_key"`
	FilterType       int     `json:"filter_type"`
	CarouselParentID string  `json:"carousel_parent_id"`
	CarouselMedia    []Item  `json:"carousel_media,omitempty"`
	User             User    `json:"user"`
	CanViewerReshare bool    `json:"can_viewer_reshare"`
	Caption          Caption `json:"caption"`
	CaptionIsEdited  bool    `json:"caption_is_edited"`
	Likes            int     `json:"like_count"`
	HasLiked         bool    `json:"has_liked"`
	// Toplikers can be `string` or `[]string`.
	Toplikers                    interface{} `json:"top_likers"`
	Likers                       []User      `json:"likers"`
	CommentLikesEnabled          bool        `json:"comment_likes_enabled"`
	CommentThreadingEnabled      bool        `json:"comment_threading_enabled"`
	HasMoreComments              bool        `json:"has_more_comments"`
	MaxNumVisiblePreviewComments int         `json:"max_num_visible_preview_comments"`
	// Previewcomments can be `string` or `[]string` or `[]Comment`.
	Previewcomments interface{} `json:"preview_comments,omitempty"`
	CommentCount    int         `json:"comment_count"`
	PhotoOfYou      bool        `json:"photo_of_you"`
	// tagged people in photo
	Tags struct {
		In []Tag `json:"in"`
	} `json:"usertags,omitempty"`
	FbUserTags           Tag    `json:"fb_user_tags"`
	CanViewerSave        bool   `json:"can_viewer_save"`
	OrganicTrackingToken string `json:"organic_tracking_token"`
	// Images contains URL images in different versions.
	Images          Images   `json:"image_versions2,omitempty"`
	OriginalWidth   int      `json:"original_width,omitempty"`
	OriginalHeight  int      `json:"original_height,omitempty"`
	ImportedTakenAt int64    `json:"imported_taken_at,omitempty"`
	Location        Location `json:"location,omitempty"`
	Lat             float64  `json:"lat,omitempty"`
	Lng             float64  `json:"lng,omitempty"`

	// Videos
	Videos            []Video `json:"video_versions,omitempty"`
	HasAudio          bool    `json:"has_audio,omitempty"`
	VideoDuration     float64 `json:"video_duration,omitempty"`
	ViewCount         float64 `json:"view_count,omitempty"`
	IsDashEligible    int     `json:"is_dash_eligible,omitempty"`
	VideoDashManifest string  `json:"video_dash_manifest,omitempty"`
	NumberOfQualities int     `json:"number_of_qualities,omitempty"`

	// Only for stories
	StoryEvents              []interface{}      `json:"story_events"`
	StoryHashtags            []interface{}      `json:"story_hashtags"`
	StoryPolls               []interface{}      `json:"story_polls"`
	StoryFeedMedia           []interface{}      `json:"story_feed_media"`
	StorySoundOn             []interface{}      `json:"story_sound_on"`
	CreativeConfig           interface{}        `json:"creative_config"`
	StoryLocations           []interface{}      `json:"story_locations"`
	StorySliders             []interface{}      `json:"story_sliders"`
	StoryQuestions           []interface{}      `json:"story_questions"`
	StoryProductItems        []interface{}      `json:"story_product_items"`
	StoryCTA                 []StoryCTA         `json:"story_cta"`
	ReelMentions             []StoryReelMention `json:"reel_mentions"`
	SupportsReelReactions    bool               `json:"supports_reel_reactions"`
	ShowOneTapFbShareTooltip bool               `json:"show_one_tap_fb_share_tooltip"`
	HasSharedToFb            int64              `json:"has_shared_to_fb"`
	Mentions                 []Mentions
	Audience                 string `json:"audience,omitempty"`
	StoryMusicStickers       []struct {
		X              float64 `json:"x"`
		Y              float64 `json:"y"`
		Z              int     `json:"z"`
		Width          float64 `json:"width"`
		Height         float64 `json:"height"`
		Rotation       float64 `json:"rotation"`
		IsPinned       int     `json:"is_pinned"`
		IsHidden       int     `json:"is_hidden"`
		IsSticker      int     `json:"is_sticker"`
		MusicAssetInfo struct {
			ID                       string `json:"id"`
			Title                    string `json:"title"`
			Subtitle                 string `json:"subtitle"`
			DisplayArtist            string `json:"display_artist"`
			CoverArtworkURI          string `json:"cover_artwork_uri"`
			CoverArtworkThumbnailURI string `json:"cover_artwork_thumbnail_uri"`
			ProgressiveDownloadURL   string `json:"progressive_download_url"`
			HighlightStartTimesInMs  []int  `json:"highlight_start_times_in_ms"`
			IsExplicit               bool   `json:"is_explicit"`
			DashManifest             string `json:"dash_manifest"`
			HasLyrics                bool   `json:"has_lyrics"`
			AudioAssetID             string `json:"audio_asset_id"`
			IgArtist                 struct {
				Pk            int    `json:"pk"`
				Username      string `json:"username"`
				FullName      string `json:"full_name"`
				IsPrivate     bool   `json:"is_private"`
				ProfilePicURL string `json:"profile_pic_url"`
				ProfilePicID  string `json:"profile_pic_id"`
				IsVerified    bool   `json:"is_verified"`
			} `json:"ig_artist"`
			PlaceholderProfilePicURL string `json:"placeholder_profile_pic_url"`
			ShouldMuteAudio          bool   `json:"should_mute_audio"`
			ShouldMuteAudioReason    string `json:"should_mute_audio_reason"`
			OverlapDurationInMs      int    `json:"overlap_duration_in_ms"`
			AudioAssetStartTimeInMs  int    `json:"audio_asset_start_time_in_ms"`
		} `json:"music_asset_info"`
	} `json:"story_music_stickers,omitempty"`
}

type Mentions struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        int64   `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	User     User    `json:"user"`
}

type StoryReelMention struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        int     `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	IsHidden int     `json:"is_hidden"`
	User     User
}

type StoryCTA struct {
	Links []struct {
		LinkType                                int         `json:"linkType"`
		WebURI                                  string      `json:"webUri"`
		AndroidClass                            string      `json:"androidClass"`
		Package                                 string      `json:"package"`
		DeeplinkURI                             string      `json:"deeplinkUri"`
		CallToActionTitle                       string      `json:"callToActionTitle"`
		RedirectURI                             interface{} `json:"redirectUri"`
		LeadGenFormID                           string      `json:"leadGenFormId"`
		IgUserID                                string      `json:"igUserId"`
		AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
	} `json:"links"`
}

type Video struct {
	Type   int    `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
	ID     string `json:"id"`
}

type Location struct {
	Pk               int64   `json:"pk"`
	Name             string  `json:"name"`
	Address          string  `json:"address"`
	City             string  `json:"city"`
	ShortName        string  `json:"short_name"`
	Lng              float64 `json:"lng"`
	Lat              float64 `json:"lat"`
	ExternalSource   string  `json:"external_source"`
	FacebookPlacesID int64   `json:"facebook_places_id"`
}

type Candidate struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type Images struct {
	Versions []Candidate `json:"candidates"`
}

type Tag struct {
	In []struct {
		User                  User        `json:"user"`
		Position              []float64   `json:"position"`
		StartTimeInVideoInSec interface{} `json:"start_time_in_video_in_sec"`
		DurationInVideoInSec  interface{} `json:"duration_in_video_in_sec"`
	} `json:"in"`
}

type Comments struct {
	item     *Item
	endpoint string
	err      error

	Items                          []Comment       `json:"comments"`
	CommentCount                   int64           `json:"comment_count"`
	Caption                        Caption         `json:"caption"`
	CaptionIsEdited                bool            `json:"caption_is_edited"`
	HasMoreComments                bool            `json:"has_more_comments"`
	HasMoreHeadloadComments        bool            `json:"has_more_headload_comments"`
	ThreadingEnabled               bool            `json:"threading_enabled"`
	MediaHeaderDisplay             string          `json:"media_header_display"`
	InitiateAtTop                  bool            `json:"initiate_at_top"`
	InsertNewCommentToTop          bool            `json:"insert_new_comment_to_top"`
	PreviewComments                []Comment       `json:"preview_comments"`
	NextMaxID                      json.RawMessage `json:"next_max_id,omitempty"`
	NextMinID                      json.RawMessage `json:"next_min_id,omitempty"`
	CommentLikesEnabled            bool            `json:"comment_likes_enabled"`
	DisplayRealtimeTypingIndicator bool            `json:"display_realtime_typing_indicator"`
	Status                         string          `json:"status"`
	//PreviewComments                []Comment `json:"preview_comments"`
}

type Caption struct {
	ID              int64  `json:"pk"`
	UserID          int64  `json:"user_id"`
	Text            string `json:"text"`
	Type            int    `json:"type"`
	CreatedAt       int64  `json:"created_at"`
	CreatedAtUtc    int64  `json:"created_at_utc"`
	ContentType     string `json:"content_type"`
	Status          string `json:"status"`
	BitFlags        int    `json:"bit_flags"`
	User            User   `json:"user"`
	DidReportAsSpam bool   `json:"did_report_as_spam"`
	MediaID         int64  `json:"media_id"`
	HasTranslation  bool   `json:"has_translation"`
}

type User struct {
	ID                         int64   `json:"pk"`
	Username                   string  `json:"username"`
	FullName                   string  `json:"full_name"`
	Biography                  string  `json:"biography"`
	ProfilePicURL              string  `json:"profile_pic_url"`
	Email                      string  `json:"email"`
	PhoneNumber                string  `json:"phone_number"`
	IsBusiness                 bool    `json:"is_business"`
	Gender                     int     `json:"gender"`
	ProfilePicID               string  `json:"profile_pic_id"`
	HasAnonymousProfilePicture bool    `json:"has_anonymous_profile_picture"`
	IsPrivate                  bool    `json:"is_private"`
	IsUnpublished              bool    `json:"is_unpublished"`
	AllowedCommenterType       string  `json:"allowed_commenter_type"`
	IsVerified                 bool    `json:"is_verified"`
	MediaCount                 int     `json:"media_count"`
	FollowerCount              int     `json:"follower_count"`
	FollowingCount             int     `json:"following_count"`
	FollowingTagCount          int     `json:"following_tag_count"`
	MutualFollowersID          []int64 `json:"profile_context_mutual_follow_ids"`
	ProfileContext             string  `json:"profile_context"`
	GeoMediaCount              int     `json:"geo_media_count"`
	ExternalURL                string  `json:"external_url"`
	HasBiographyTranslation    bool    `json:"has_biography_translation"`
	ExternalLynxURL            string  `json:"external_lynx_url"`
	BiographyWithEntities      struct {
		RawText  string        `json:"raw_text"`
		Entities []interface{} `json:"entities"`
	} `json:"biography_with_entities"`
	UsertagsCount                int        `json:"usertags_count"`
	HasChaining                  bool       `json:"has_chaining"`
	IsFavorite                   bool       `json:"is_favorite"`
	IsFavoriteForStories         bool       `json:"is_favorite_for_stories"`
	IsFavoriteForHighlights      bool       `json:"is_favorite_for_highlights"`
	CanBeReportedAsFraud         bool       `json:"can_be_reported_as_fraud"`
	ShowShoppableFeed            bool       `json:"show_shoppable_feed"`
	ShoppablePostsCount          int        `json:"shoppable_posts_count"`
	ReelAutoArchive              string     `json:"reel_auto_archive"`
	HasHighlightReels            bool       `json:"has_highlight_reels"`
	PublicEmail                  string     `json:"public_email"`
	PublicPhoneNumber            string     `json:"public_phone_number"`
	PublicPhoneCountryCode       string     `json:"public_phone_country_code"`
	ContactPhoneNumber           string     `json:"contact_phone_number"`
	CityID                       int64      `json:"city_id"`
	CityName                     string     `json:"city_name"`
	AddressStreet                string     `json:"address_street"`
	DirectMessaging              string     `json:"direct_messaging"`
	Latitude                     float64    `json:"latitude"`
	Longitude                    float64    `json:"longitude"`
	Category                     string     `json:"category"`
	BusinessContactMethod        string     `json:"business_contact_method"`
	IncludeDirectBlacklistStatus bool       `json:"include_direct_blacklist_status"`
	Byline                       string     `json:"byline"`
	SocialContext                string     `json:"social_context,omitempty"`
	SearchSocialContext          string     `json:"search_social_context,omitempty"`
	MutualFollowersCount         float64    `json:"mutual_followers_count"`
	LatestReelMedia              int64      `json:"latest_reel_media,omitempty"`
	IsCallToActionEnabled        bool       `json:"is_call_to_action_enabled"`
	FbPageCallToActionID         string     `json:"fb_page_call_to_action_id"`
	Zip                          string     `json:"zip"`
	Friendship                   Friendship `json:"friendship_status"`
}

type Friendship struct {
	IncomingRequest bool `json:"incoming_request"`
	FollowedBy      bool `json:"followed_by"`
	OutgoingRequest bool `json:"outgoing_request"`
	Following       bool `json:"following"`
	Blocking        bool `json:"blocking"`
	IsPrivate       bool `json:"is_private"`
	Muting          bool `json:"muting"`
	IsMutingReel    bool `json:"is_muting_reel"`
}

type Comment struct {
	idstr string

	ID                             int64     `json:"pk"`
	Text                           string    `json:"text"`
	Type                           int       `json:"type"`
	User                           User      `json:"user"`
	UserID                         int64     `json:"user_id"`
	BitFlags                       int       `json:"bit_flags"`
	ChildCommentCount              int       `json:"child_comment_count"`
	CommentIndex                   int       `json:"comment_index"`
	CommentLikeCount               int       `json:"comment_like_count"`
	ContentType                    string    `json:"content_type"`
	CreatedAt                      int64     `json:"created_at"`
	CreatedAtUtc                   int64     `json:"created_at_utc"`
	DidReportAsSpam                bool      `json:"did_report_as_spam"`
	HasLikedComment                bool      `json:"has_liked_comment"`
	InlineComposerDisplayCondition string    `json:"inline_composer_display_condition"`
	OtherPreviewUsers              []User    `json:"other_preview_users"`
	PreviewChildComments           []Comment `json:"preview_child_comments"`
	NextMaxChildCursor             string    `json:"next_max_child_cursor,omitempty"`
	HasMoreTailChildComments       bool      `json:"has_more_tail_child_comments,omitempty"`
	NextMinChildCursor             string    `json:"next_min_child_cursor,omitempty"`
	HasMoreHeadChildComments       bool      `json:"has_more_head_child_comments,omitempty"`
	NumTailChildComments           int       `json:"num_tail_child_comments,omitempty"`
	NumHeadChildComments           int       `json:"num_head_child_comments,omitempty"`
	Status                         string    `json:"status"`
}