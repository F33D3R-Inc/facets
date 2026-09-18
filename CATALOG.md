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
| **FeedShell (vertical)** | full-screen one-post-at-a-time vertical pager (TikTok / Shorts / Reels) with snap scrolling and an action column | ✅ `ui/feedshell.fct` (proven via `smoke_new2.fct`) |
| TopBar | sticky blurred header with back + title + sub | ✅ |
| **BottomNav** | 5-tab mobile bar (Home · Discover · + · Inbox · Me) | ✅ `ui/bottomnav.fct` (proven via `smoke_new.fct`) |
| Pill / buttons | primary, light, ghost | ✅ (classes) |
| SignUpCard · GuestBanner · Login | the signed-out surfaces | ✅ |
| **Popover / Menu** | anchored menu (the "More" list, a post's "…") | ✅ Menu ✅ `ui/menu.fct` (a modal sheet, not an anchored dropdown — proven via `smoke.fct`) · Popover ✅ `ui/popover.fct` |
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
| **ShortsCard** | full-bleed vertical video with the right-hand action column (like, comment, share, sound) and bottom caption/author | ✅ `media/shortscard.fct` (proven via `smoke_new2.fct`) |
| **VideoPage (YouTube)** | player · title · channel row with Subscribe · like/dislike/share/save bar · description fold · comments · Up-next rail | ✅ `media/videopage.fct` (proven via `smoke_new3.fct`; comments reuse `content/comments.fct` at the call site — no dedicated glyph for "dislike" in `ui/icon.fct` yet, so that button is plain text) |
| **VideoTile** | thumbnail + duration badge + title + channel + views · time | ✅ `media/videotile.fct` (proven via `smoke_new.fct`) |
| **Playlist / Series** | ordered list of VideoTiles with progress | ✅ `media/playlist.fct` (proven via `smoke_new2.fct`) |
| **LinkPreview** | og-card for a URL in a post | 🟡 `content/linkpreview.fct` (proven via `smoke_new4.fct`; fetch/extract via two `proc`s using `httpGet` + `split`, render side takes the resolved fields — never clickable, since `link "…" -> "https://…"` is refused and there is no external destination in the language yet, a real gap this file names rather than fakes) |
| **Poll** | options with vote bars, one vote per actor (`exists`) | ✅ `social/poll.fct` (proven via `smoke_new2.fct`) |
| **Hashtag / Mention** | `#tag` `@handle` autolinks + `/tag/:tag` page | ✅ |
| **Repost / Quote** | quote via `Tweet(t.quoted)` in the PostCard slot | ✅ |
| **Bookmarks · Lists · Communities** | saved posts, curated user lists, topic groups | ✅ `social/bookmarkslist.fct` · `graph/userlist.fct` · `social/community.fct` (proven via `smoke_new3.fct`) |
| **Search / Explore** | SearchBox + trending + results tabs (Top · Latest · People · Media) | 🟡 (SearchBox, Trends) |
| **Notifications** | NotificationItem list + UnreadBadge; fan-out actions | ✅ atoms + page `notify/notificationspage.fct` (proven via `smoke_new3.fct`) |

## 2 · Going live (`live/`) — one set for Twitch, Kick, YouTube Live, TikTok LIVE, cam sites

| Facet | What it is | Status |
|---|---|---|
| **LivePlayer** | the stream frame: `video` (HLS URL) + LIVE pill + viewer count + uptime | ✅ `live/liveplayer.fct` |
| **StreamInfo** | below the player: channel row, Follow, Subscribe, Tip, category + tags | ✅ `live/streaminfo.fct` |
| **ChatPanel** | realtime message list over SSE, pinned message, guest gate, composer | ✅ `live/chatpanel.fct` (emotes ✅ — via composition: a `ToggleButtonQuiet` + `EmotePicker` beside the call site's own `ChatPanel`, `insertEmote` appending to the same `draft` cell `ChatPanel` binds, proven via `smoke_new5.fct`; slow mode ⬜ — timer) |
| **ChatMessage** | role glyph (broadcaster · mod · sub) + coloured name + richtext body + mod tools | ✅ `live/chatmessage.fct` |
| **Emote / Sticker picker** | grid in a Sheet; inserts a token | ✅ `live/emotepicker.fct` (proven via `smoke_new3.fct`; the token lands via a fixed `insertEmote(token)` action, same convention as `Poll`'s `vote`) |
| **TipButton + TipSheet** | preset amounts, note, calls `tip(stream, cents, note)` | ✅ `live/tipsheet.fct` (custom amount ⬜) |
| **GoalBar** | "Tip goal" progress with label — one facet also serves sub goals and fundraisers | ✅ `live/goalbar.fct` |
| **TopTippers / Leaderboard** | ranked list from `sum(t.cents in Tip where …)` per user | ✅ `live/leaderboard.fct` (proven via `smoke_new3.fct`; still no group-by — the host supplies rank via `for` over its own roster with a filtered `sum` each, as before) |
| **Alerts / RecentTips** | the latest tips, live | ✅ `live/recenttips.fct` (timed banner ⬜ — timer) |
| **SubscribeButton** | calls `subscribe(stream)` / `unsubscribe` | ✅ `live/subscribebutton.fct` (tiers ✅ — `SubscribeTier`/`SubscribeTierList` call the host's own `subscribeTier(stream, tier)`/`unsubscribeTier(stream)`, a separate pair of actions rather than widening `subscribe`'s arity, proven via `smoke_new5.fct`) |
| **GiftPicker** | virtual gifts grid (TikTok/Kick) — TipSheet with pictures | ✅ `live/giftpicker.fct` (proven via `smoke_new5.fct`; each gift is a row — picture, name, price, then the `button` that calls `tip(stream, cents, name)` — since a `button` renders no child nodes, so the picture cannot sit inside it; the picture is a generated image seeded by the gift's name, the same technique `ui/Avatar`/`MyQRCode` already use, not an `ui/icon.fct` glyph, since "gift"/"crown"/"coin" are not in that set) |
| **ViewerList** | who is watching, with roles | ✅ `live/viewerlist.fct` (proven via `smoke_new2.fct`) |
| **ModTools** | delete / ban per message — `requires moderator(stream)` | ✅ in ChatMessage (slow mode ⬜) |
| **StreamCard** | live tile: thumbnail + LIVE + viewers + title + channel + category | ✅ `live/streamcard.fct` |
| **BrowseGrid** | StreamCard grid, most watched first, optional category | ✅ `live/browsegrid.fct` |
| **CategoryChip / Tag** | one pill → `/browse/:category` | ✅ `live/categorychip.fct` |
| **Schedule** | upcoming streams calendar | ✅ `live/schedule.fct` (proven via `smoke_new2.fct`) |
| **VOD list / Clips** | past broadcasts = VideoTiles; Clips = ShortsCard | ✅ by reuse |
| **GoLivePanel (creator)** | title, category, masked stream key (`@secret`, reveal toggle), go live / end | ✅ `live/golivepanel.fct` (thumbnail upload ✅ — the same `UploadField` scaffold `ProfileEdit` uses, `updateStream(id, title, category, thumbnail)` widened to write `Stream(id).thumb`, proven via `smoke_new5.fct`) |
| **Creator dashboard** | viewers over time, income, top clips — StatChips + a chart | ⬜ (chart — needs a `chart` node or SVG facet) |
| **AgeGate / ContentWarning** | interstitial requiring confirmation, remembered per browser | ✅ `live/agegate.fct` |
| **PrivateShow / Paywall** | gated region: `requires subscriber(channel)` around the player | ✅ `live/paywall.fct` (proven via `smoke_new3.fct`; same gated/ungated `if` split as `AgeGate` — the render side, alongside a `requires subscriber(channel)` on whatever action the gated content itself calls) |

## 3 · Messaging & social graph (`chat/`, `graph/`)

| Facet | What it is | Status |
|---|---|---|
| **DMThread** | 1:1 messages, live over SSE, read receipts ("Seen"), `@e2e` sealed bodies when the host declares it | ✅ `chat/dmthread.fct` (groups ⬜ — need `let id = add …` to seed members; typing ⬜ — timer) |
| **InboxList** | conversations with unread badges and last-message preview | ✅ `chat/inboxlist.fct` |
| **ContactCard / Friend request** | add by handle, Requested / Message states, accept / decline | ✅ `graph/contactcard.fct`, `graph/friendrequestrow.fct` (QR ✅ — `ContactQR`, the same generated-image technique `commerce/PayButton`'s `MyQRCode` already proves, applied to an "add me" payload instead of a payment one; proven via `smoke_new5.fct`; scanning it back is the same camera/lang gap already named for `PayButton`) |
| FollowButton · WhoToFollow | | ✅ |
| **BlockMuteMenu** | block, mute, report — a sheet until there is a popover | ✅ `graph/usermenu.fct` |
| **ReportSheet** | reasons list → `report(target, reason)` | ✅ `graph/reportsheet.fct` |

## 4 · Super-app (WeChat)

| Facet | What it is | Status |
|---|---|---|
| **Moments** | a friends-only feed of PostCards with a cover image and a 9-grid Gallery | ✅ `social/moments.fct` (proven via `smoke_new4.fct`; friends-only is the host's `exists` in its own `for`, same split `social/CommunityGrid` uses for policy) |
| **Wallet** | balance, top-up, transfer, red packet — on the billing ledger | ✅ `commerce/wallet.fct` (proven via `smoke_new4.fct`; a `WalletTx` ledger + `sum` for the balance, same shape `live/Leaderboard` already settles a ledger with; custom amount ⬜ — same text-to-money gap `live/TipSheet` already names) |
| **PayButton / QRPay** | pay a merchant, show my QR, scan | 🟡 `commerce/paybutton.fct` (`PayButton` + `MyQRCode`, proven via `smoke_new4.fct`; scan ⬜ — camera, lang) |
| **MiniProgramTile / Drawer** | grid of apps; a mini-program is another `.fct` mounted at a route | ✅ `ui/miniprogramdrawer.fct` (proven via `smoke_new4.fct`; the tile is a `Link` into whatever route a host mounts the program at, a `Sheet` is the drawer at the call site) |
| **OfficialAccount** | a channel page: articles (richtext) + subscribe + menu | ✅ `content/officialaccount.fct`'s `OfficialAccountHeader` (proven via `smoke_new4.fct`; articles = `content/Article` (already ✅), menu = `navigation/TabBar` (already ✅), subscribe = `live/SubscribeButton` reused whole — the header for a brand rather than a person was the only new piece) |
| **Channels (video)** | = ShortsCard feed | ✅ by reuse |
| **Sticker keyboard** | = Emote picker | ✅ by reuse |
| **Location / Nearby** | map card | ✅ `ui/locationcard.fct` (`LocationCard` + `NearbyList`/`NearbyRow`, proven via `smoke_new4.fct`; the map is a static image the host builds the URL for — no live/pannable map, no `embed` node in the language — and "nearby" is the host's own already-computed distance, not live device geolocation) |

## 5 · Accounts & commerce (cross-cutting)

| Facet | What it is | Status |
|---|---|---|
| Login / SignUpCard | | ✅ |
| **SocialLogin buttons** | Google / Apple | 🟡 (one OIDC provider — lang for a second) |
| **SettingsPage** | sections list + toggles (`toggle` node) | ✅ `profile/settingspage.fct` (proven via `smoke_new.fct`) |
| **ProfileEdit** | avatar/banner `upload`, bio, links | ✅ `profile/profileedit.fct` (proven via `smoke_new2.fct`) |
| **Checkout / Pricing** | plans grid, PayButton | ✅ `commerce/checkout.fct` (`CheckoutSummary`/`CheckoutLineItem`/`CheckoutTotal`/`CheckoutActions`, proven via `smoke_new4.fct`; the grid was already `marketing/PricingTable`, the button is `commerce/PayButton` — this file is the itemised summary between them) |
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
9. **An external `link` destination** — the README already names this ("nothing here can leave the site": `link "…" -> "https://…"` is a refused path, and an absolute-URL `Link` renders as inert text). It surfaced again this round: `content/LinkPreview`'s og-card can fetch and render a shared URL's title/description/image, but can never make the card itself clickable — the one thing every platform's link preview is *for*. A repo-links footer and an og-card are the same missing piece.
10. **Text-to-money conversion** — an `input bind` cell is always `text`; there is no builtin to turn "12.50" into a `money`/`int` value. `live/TipSheet`'s "custom amount" gap and `commerce/Wallet`'s three sheets (top-up, transfer, red packet) all stop at presets for exactly this reason.

Everything else above is library work on the language as it stands today.
