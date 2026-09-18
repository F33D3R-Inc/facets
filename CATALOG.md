# The facet catalog — what a social platform is made of

The target: enough facets that a **creator-social platform with live streaming**
(the union of X, TikTok, YouTube, Twitch, Kick, and a cam site) and a
**super-app** (WeChat: chat, moments, pay, mini-programs) are both *assembled*
from this library, not written. One facet per need, never one per platform — a
Twitch "channel", a YouTube "channel", a TikTok "profile" and an X "profile" are
one `ProfileHeader` with different data in it.

Status: ✅ shipped in this repo · 🟡 partial · ⬜ to build. Every facet below is
expressible in the language today unless marked **needs lang**, and those are
listed at the end. Bold rows are the ones that unblock the most screens.

## 0 · Foundations (`ui/`)

| Facet | What it is | Status |
|---|---|---|
| **Look** | the theme + `x-` layout vocabulary every atom uses | ✅ |
| **Icons** | 30 glyph masks (`icon "heart"`, `class "x-glyph-heart"`) | ✅ (grow to ~60: gift, coin, crown, mic-off, camera, screen, pin-map, shield, flag, translate, poll, emoji, sticker, image, wallet, qr, scan, moments, mini-app) |
| AppShell (3-col) | rail · main · aside, sticky | ✅ (`wireframes/shell.fct`, on `layout/vocabulary.fct`'s `x-l-app3`) |
| **TheaterShell** | the live layout: video left, chat right, no rail; collapses to video-over-chat on mobile | ✅ `ui/theatershell.fct` (proven via `smoke_new.fct`) |
| **FeedShell (vertical)** | full-screen one-post-at-a-time vertical pager (TikTok / Shorts / Reels) with snap scrolling and an action column | ⬜ |
| TopBar | sticky blurred header with back + title + sub | ✅ |
| **BottomNav** | 5-tab mobile bar (Home · Discover · + · Inbox · Me) | ✅ `ui/bottomnav.fct` (proven via `smoke_new.fct`) |
| Pill / buttons | primary, light, ghost | ✅ (classes) |
| SignUpCard · GuestBanner · Login | the signed-out surfaces | ✅ |
| **Popover / Menu** | anchored menu (the "More" list, a post's "…") | 🟡 Menu ✅ `ui/menu.fct` (a modal sheet, not an anchored dropdown — proven via `smoke.fct`) · Popover ⬜ (`overlay` is centered; needs a `popover` node or an anchored class — see lang) |
| **Sheet** | bottom sheet (mobile dialogs: share, gift picker, comments) | ✅ `ui/sheet.fct` (proven via `smoke_new.fct`) |
| Toast | transient confirmation ("Copied", "Followed") | ✅ `ui/toast.fct` (proven via `smoke.fct`; auto-dismiss still needs a timer — lang) |
| Skeleton / Empty / Error states | one facet each, used everywhere | ✅ `ui/skeleton.fct` (SkeletonLine/SkeletonPost) · `ui/emptystate.fct` · `ui/errorstate.fct` (all proven via `smoke.fct`) |
| Tabs · Badge · Avatar · VerifiedBadge · AuthorRow · UserChip | identity atoms | ✅ |
| **StatChip** | "240.1M Followers" — `compact` number + label | ✅ `ui/statchip.fct` (proven via `smoke_new.fct`) |
| RelativeTime | `ago(ts)` in a muted span with a full-date title | ✅ `ui/relativetime.fct` (no tooltip — no generic `title` attribute in the language yet; proven via `smoke_new.fct`) |

## 1 · Content (`social/`, `media/`)

| Facet | What it is | Status |
|---|---|---|
| PostCard (slot) · QuoteCard · EngagementBar · MediaCard | the post | ✅ |
| **Composer** | text + attach (image/video) + poll + audience + schedule; shows `pending`/`failed` | 🟡 (`ComposeBox` is text-only; `upload` exists) |
| **ThreadView** | a post with its replies, nested one level, inline reply composer | ✅ `social/threadview.fct` (proven via `smoke_new.fct`) |
| **Gallery** | 1–4 images in the X grid; tap → lightbox (`overlay`) | ✅ `ui/gallery.fct` (proven via `smoke.fct`; lightbox composition is a call-site `overlay`) |
| **VideoPlayer** | `video` + poster + custom controls + mute toggle; autoplay muted in feeds | 🟡 (`MediaCard`; controls are native) |
| **ShortsCard** | full-bleed vertical video with the right-hand action column (like, comment, share, sound) and bottom caption/author | ⬜ |
| **VideoPage (YouTube)** | player · title · channel row with Subscribe · like/dislike/share/save bar · description fold · comments · Up-next rail | ⬜ (all atoms exist except Up-next = a `for` of `VideoTile`) |
| **VideoTile** | thumbnail + duration badge + title + channel + views · time | ✅ `media/videotile.fct` (proven via `smoke_new.fct`) |
| **Playlist / Series** | ordered list of VideoTiles with progress | ⬜ |
| **LinkPreview** | og-card for a URL in a post | ⬜ (needs a service `call` to unfurl — lang: service exists) |
| **Poll** | options with vote bars, one vote per actor (`exists`) | ⬜ |
| **Hashtag / Mention** | `#tag` `@handle` autolinks + `/tag/:tag` page | ✅ |
| **Repost / Quote** | quote via `Tweet(t.quoted)` in the PostCard slot | ✅ |
| **Bookmarks · Lists · Communities** | saved posts, curated user lists, topic groups | ⬜ (entities + FollowButton-style toggles) |
| **Search / Explore** | SearchBox + trending + results tabs (Top · Latest · People · Media) | 🟡 (SearchBox, Trends) |
| **Notifications** | NotificationItem list + UnreadBadge; fan-out actions | ✅ atoms / ⬜ page |

## 2 · Going live (`live/`) — one set for Twitch, Kick, YouTube Live, TikTok LIVE, cam sites

| Facet | What it is | Status |
|---|---|---|
| **LivePlayer** | the stream frame: `video` (HLS URL) + LIVE pill + viewer count + uptime | ✅ `live/liveplayer.fct` |
| **StreamInfo** | below the player: channel row, Follow, Subscribe, Tip, category + tags | ✅ `live/streaminfo.fct` |
| **ChatPanel** | realtime message list over SSE, pinned message, guest gate, composer | ✅ `live/chatpanel.fct` (emotes, slow mode ⬜) |
| **ChatMessage** | role glyph (broadcaster · mod · sub) + coloured name + richtext body + mod tools | ✅ `live/chatmessage.fct` |
| **Emote / Sticker picker** | grid in a Sheet; inserts a token | ⬜ |
| **TipButton + TipSheet** | preset amounts, note, calls `tip(stream, cents, note)` | ✅ `live/tipsheet.fct` (custom amount ⬜) |
| **GoalBar** | "Tip goal" progress with label — one facet also serves sub goals and fundraisers | ✅ `live/goalbar.fct` |
| **TopTippers / Leaderboard** | ranked list from `sum(t.cents in Tip where …)` per user | ⬜ (needs group-by; today: `for` over Users with a filtered `sum` each) |
| **Alerts / RecentTips** | the latest tips, live | ✅ `live/recenttips.fct` (timed banner ⬜ — timer) |
| **SubscribeButton** | calls `subscribe(stream)` / `unsubscribe` | ✅ `live/subscribebutton.fct` (tiers ⬜) |
| **GiftPicker** | virtual gifts grid (TikTok/Kick) — TipSheet with pictures | 🟡 (TipSheet; images ⬜) |
| **ViewerList** | who is watching, with roles | ⬜ |
| **ModTools** | delete / ban per message — `requires moderator(stream)` | ✅ in ChatMessage (slow mode ⬜) |
| **StreamCard** | live tile: thumbnail + LIVE + viewers + title + channel + category | ✅ `live/streamcard.fct` |
| **BrowseGrid** | StreamCard grid, most watched first, optional category | ✅ `live/browsegrid.fct` |
| **CategoryChip / Tag** | one pill → `/browse/:category` | ✅ `live/categorychip.fct` |
| **Schedule** | upcoming streams calendar | ⬜ |
| **VOD list / Clips** | past broadcasts = VideoTiles; Clips = ShortsCard | ✅ by reuse |
| **GoLivePanel (creator)** | title, category, masked stream key (`@secret`, reveal toggle), go live / end | ✅ `live/golivepanel.fct` (thumbnail upload ⬜) |
| **Creator dashboard** | viewers over time, income, top clips — StatChips + a chart | ⬜ (chart — needs a `chart` node or SVG facet) |
| **AgeGate / ContentWarning** | interstitial requiring confirmation, remembered per browser | ✅ `live/agegate.fct` |
| **PrivateShow / Paywall** | gated region: `requires subscriber(channel)` around the player | ⬜ (policy + region guard exist) |

## 3 · Messaging & social graph (`chat/`, `graph/`)

| Facet | What it is | Status |
|---|---|---|
| **DMThread** | 1:1 messages, live over SSE, read receipts ("Seen"), `@e2e` sealed bodies when the host declares it | ✅ `chat/dmthread.fct` (groups ⬜ — need `let id = add …` to seed members; typing ⬜ — timer) |
| **InboxList** | conversations with unread badges and last-message preview | ✅ `chat/inboxlist.fct` |
| **ContactCard / Friend request** | add by handle, Requested / Message states, accept / decline | ✅ `graph/contactcard.fct`, `graph/friendrequestrow.fct` (QR ⬜) |
| FollowButton · WhoToFollow | | ✅ |
| **BlockMuteMenu** | block, mute, report — a sheet until there is a popover | ✅ `graph/usermenu.fct` |
| **ReportSheet** | reasons list → `report(target, reason)` | ✅ `graph/reportsheet.fct` |

## 4 · Super-app (WeChat)

| Facet | What it is | Status |
|---|---|---|
| **Moments** | a friends-only feed of PostCards with a cover image and a 9-grid Gallery | ⬜ (PostCard + Gallery + policy) |
| **Wallet** | balance, top-up, transfer, red packet — on the billing ledger | ⬜ |
| **PayButton / QRPay** | pay a merchant, show my QR, scan | ⬜ (QR = `image` from a QR service; scan = camera — lang: no camera) |
| **MiniProgramTile / Drawer** | grid of apps; a mini-program is another `.fct` mounted at a route | ⬜ (imports + routes exist) |
| **OfficialAccount** | a channel page: articles (richtext) + subscribe + menu | ⬜ (richtext + tabs) |
| **Channels (video)** | = ShortsCard feed | ✅ by reuse |
| **Sticker keyboard** | = Emote picker | ✅ by reuse |
| **Location / Nearby** | map card | ⬜ (needs an embed node or an `image` static map) |

## 5 · Accounts & commerce (cross-cutting)

| Facet | What it is | Status |
|---|---|---|
| Login / SignUpCard | | ✅ |
| **SocialLogin buttons** | Google / Apple | 🟡 (one OIDC provider — lang for a second) |
| **SettingsPage** | sections list + toggles (`toggle` node) | ✅ `profile/settingspage.fct` (proven via `smoke_new.fct`) |
| **ProfileEdit** | avatar/banner `upload`, bio, links | ⬜ |
| **Checkout / Pricing** | plans grid, PayButton | ⬜ (billing exists) |
| **AdminTable** | auto-admin already exists in the runtime | ✅ |

## What still needs the language (small, ranked)

1. **`popover bind cell:`** (anchored overlay) — menus, the "…" on every card, the gift picker. One node lowered like `overlay`; positioned beside the previous sibling.
2. **A timer** — `after 5s: cell = false` (toast, alert, slow-mode countdown). Client-only, one statement kind; placement is trivially client.
3. **`on open -> action`** on a view — view counts, "mark seen", watch history.
4. **Second OIDC provider** — Apple beside Google.
5. **Group-by aggregates** — leaderboards without a `for` over every user.
6. **A `chart` node** (bar/line from a `for`) — creator dashboards. Could be an SVG-emitting facet if `svg`/`style` per element were allowed; a node is cleaner.
7. **Camera/QR scan** — a `scan bind cell` control (getUserMedia) — the one WeChat surface no facet can fake.
8. **`let id = add Entity {…}`** — bind the id of a row an action just added, so one action can create a conversation and its member rows. Without it a group chat has no way to seed membership; 1:1 threads work around it by carrying both participants on the row.

Everything else above is library work on the language as it stands today.
