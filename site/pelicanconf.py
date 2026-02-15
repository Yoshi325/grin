AUTHOR = "Charles Y."
SITENAME = "grin"
SITEURL = ""

PATH = "content"
TIMEZONE = "America/New_York"
DEFAULT_LANG = "en"

# Theme
THEME = "themes/papyrus"

# Plugins
PLUGINS = ["readtime", "neighbors", "pelican-toc"]
PLUGIN_PATHS = ["pelican-plugins"]

# Static paths
STATIC_PATHS = ["images"]

# URL settings
PAGE_URL = "{slug}/"
PAGE_SAVE_AS = "{slug}/index.html"

# Disable blog features (this is a project site, not a blog)
DIRECT_TEMPLATES = ["index"]
ARTICLE_PATHS = []
DISPLAY_PAGES_ON_MENU = True

# Template overrides for custom landing page
THEME_TEMPLATES_OVERRIDES = ["templates"]

# Feed (disable for local dev)
FEED_ALL_ATOM = None
CATEGORY_FEED_ATOM = None
TRANSLATION_FEED_ATOM = None
AUTHOR_FEED_ATOM = None
AUTHOR_FEED_RSS = None
