# Linkor Usage Guide

## Configuration

Linkor reads configuration from environment variables. You can also place a
`.env` file in the working directory — it will be loaded automatically at
startup. Environment variables set in the shell take precedence over `.env`
values.

Copy the example file to get started:

```bash
cp .env.example .env
```

See the [README](../README.md#configuration) for the full list of variables.

## Admin Panel

Access the admin panel at:

```
http://your-server:8080/admin/management?key=YOUR_ADMIN_SECRET_KEY
```

Replace `YOUR_ADMIN_SECRET_KEY` with the value of the `ADMIN_SECRET_KEY` environment
variable you set when starting the server.

### Managing Links

#### Creating a Link

1. Click **+ New Link** in the top-right corner.
2. Fill in the fields:
   - **Name** (required): A display name for the link (e.g. "Product Landing Page").
   - **Short URL** (optional): A custom slug for the short link. Leave empty to
     auto-generate a random 6-character code. Only letters, digits, hyphens, and
     underscores are allowed (2-64 characters).
   - **Target URL** (required): The full destination URL to redirect to.
3. Click **Save**.

#### Editing a Link

1. Click **Edit** on the link row.
2. Modify any field and click **Save**.

Note: Changing the short URL will break existing links pointing to the old slug.

#### Deleting a Link

1. Click **Delete** on the link row.
2. Confirm the deletion in the dialog.

Deletion removes all associated click tracking data.

### Statistics

The header bar shows aggregated stats across all links:

- **Links**: Total number of short links.
- **Total Clicks**: Sum of all clicks across all links.
- **Unique Clicks**: Sum of all unique clicks across all links.

Per-link statistics (Clicks and Unique columns) are visible in the table.

## Short Link Redirects

Short links are accessed at the root path:

```
http://your-server:8080/SHORT_URL
```

For example, if you created a link with the short URL `docs`, visitors to
`http://your-server:8080/docs` will be redirected (HTTP 302) to the target URL.

### Click Tracking

- **Clicks**: Every visit to a short link increments the click counter.
- **Unique Clicks**: Tracked by hashing the visitor's IP address (SHA-256). Each
  unique IP is counted once per link, regardless of how many times they visit.
  The system respects `X-Forwarded-For` and `X-Real-IP` headers for clients
  behind proxies.

## API Reference

All API endpoints are protected and require the admin key. Send it as the
`X-Admin-Key` header with every request. Requests without a valid key receive a
`401 Unauthorized` response.

### List Links

```
GET /admin/api/links
```

Returns a JSON array of all links.

### Create Link

```
POST /admin/api/links
Content-Type: application/json

{
  "name": "My Link",
  "short_url": "custom-slug",
  "target_url": "https://example.com"
}
```

`short_url` is optional; omit or send empty string to auto-generate.

### Update Link

```
PUT /admin/api/links/{id}
Content-Type: application/json

{
  "name": "Updated Name",
  "short_url": "new-slug",
  "target_url": "https://new-target.com"
}
```

### Delete Link

```
DELETE /admin/api/links/{id}
```

Returns `{"status": "deleted"}` on success.
