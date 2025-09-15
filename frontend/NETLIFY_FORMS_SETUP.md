# Netlify Forms Setup Guide

This guide explains how to set up and manage the Netlify Forms implementation for Cloud Partner Pro, which has been migrated from EmailJS to provide a more reliable and professional email system.

## 📋 Overview

The project now uses **Netlify Forms** with custom email templates for:
1. **Contact Form** - General inquiries and service requests
2. **Quote Request Form** - Detailed pricing calculator submissions

## 🚀 Quick Setup

### 1. Netlify Configuration
The site is configured with:
- Form detection enabled (automatic when forms are properly structured)
- Custom email templates configured via Netlify Dashboard
- Professional email routing to `info@cloudpartner.pro`

**Important:** Email templates MUST be configured through the Netlify Dashboard GUI, not via `netlify.toml`.

### 2. Forms Implementation
Both forms are now implemented with:
- ✅ Netlify form attributes (`data-netlify="true"`)
- ✅ Honeypot spam protection
- ✅ Proper form naming for email routing
- ✅ Static HTML versions for build detection

## 📧 Email Templates

### Dashboard Configuration Required
**Important:** Email templates must be configured through the Netlify Dashboard:
1. Go to your site's Netlify Dashboard
2. Navigate to Forms → Form notifications
3. Select the form (contact-form or quote-request-form)
4. Configure email notifications with custom templates

### Contact Form Template
**Template Files:** `netlify/emails/contact-form.html` and `contact-form.txt`
**Configure in Dashboard for:** `contact-form`

**Features:**
- Professional blue color scheme
- Clear contact information display
- Services of interest highlighting
- Mobile-responsive design
- Priority badge for new leads

**Fields Available:**
- `{{ name }}` - Full name
- `{{ email }}` - Email address
- `{{ company }}` - Company name
- `{{ phone }}` - Phone number
- `{{ services }}` - Selected services
- `{{ message }}` - Message content

### Quote Request Template
**Template Files:** `netlify/emails/quote-request-form.html` and `quote-request-form.txt`
**Configure in Dashboard for:** `quote-request-form`

**Features:**
- Professional green color scheme (differentiating from contact)
- Detailed pricing breakdown display
- Service requirements summary
- Quote badge for urgent attention
- Mobile-responsive design

**Fields Available:**
- `{{ name }}`, `{{ email }}`, `{{ phone }}`, `{{ company }}` - Contact info
- `{{ serviceType }}`, `{{ complexity }}`, `{{ count }}` - Service details
- `{{ basePrice }}`, `{{ variablePrice }}`, `{{ complexityMultiplier }}`, `{{ totalEstimate }}` - Pricing
- `{{ requirements }}` - Additional requirements

## 🔧 Form Configuration

### Contact Form
```jsx
<Form data-netlify="true" data-netlify-honeypot="bot-field" method="POST" name="contact-form">
  <input type="hidden" name="form-name" value="contact-form" />
  <div style={{ display: 'none' }}>
    <label>Don't fill this out if you're human: <input name="bot-field" /></label>
  </div>
  {/* form fields */}
</Form>
```

### Quote Request Form
```jsx
<form data-netlify="true" data-netlify-honeypot="bot-field" method="POST" name="quote-request-form">
  <input type="hidden" name="form-name" value="quote-request-form" />
  <input name="bot-field" />
  {/* form fields */}
</form>
```

## 🎨 Email Template Configuration

### Setting Up Templates in Netlify Dashboard

**Step-by-Step Configuration:**

1. **Deploy Your Site First**
   - Push code to your repository
   - Let Netlify build and deploy the site
   - Forms will be automatically detected

2. **Access Form Settings**
   - Go to Netlify Dashboard → Your Site
   - Navigate to **Forms** tab
   - You should see both forms listed:
     - `contact-form`
     - `quote-request-form`

3. **Configure Contact Form Email**
   - Click on `contact-form`
   - Go to **Settings & usage** → **Form notifications**
   - Click **Add notification** → **Email notification**
   - Set **Email to notify**: `info@cloudpartner.pro`
   - Set **Subject line**: `New Contact Form Submission - Cloud Partner Pro`
   - **Custom email template**: Copy HTML from `EMAIL_TEMPLATES.md` → Contact Form section
   - **Save notification**

4. **Configure Quote Request Form Email**
   - Click on `quote-request-form`
   - Go to **Settings & usage** → **Form notifications**
   - Click **Add notification** → **Email notification**
   - Set **Email to notify**: `info@cloudpartner.pro`
   - Set **Subject line**: `New Quote Request - Cloud Partner Pro`
   - **Custom email template**: Copy HTML from `EMAIL_TEMPLATES.md` → Quote Request Form section
   - **Save notification**

### Template Files Reference
**See `EMAIL_TEMPLATES.md` for complete template code to copy into Netlify Dashboard**

Template configurations needed:
- **Contact Form:** Use templates from `EMAIL_TEMPLATES.md` → Contact Form section
- **Quote Form:** Use templates from `EMAIL_TEMPLATES.md` → Quote Request Form section

### Brand Colors
- **Contact Form:** Blue theme (`#1e40af`, `#3b82f6`)
- **Quote Form:** Green theme (`#059669`, `#10b981`)

### Customizing Templates
1. Edit HTML templates in `EMAIL_TEMPLATES.md`
2. Copy updated content to Netlify Dashboard
3. Test with form submissions
4. Modify colors, fonts, layout as needed

### Adding New Fields
1. Add field to React form component
2. Update static HTML version in `/public/`
3. Add field to email template with `{{ field_name }}` syntax
4. Update template in Netlify Dashboard
5. Deploy changes

## 🧪 Testing

### Local Testing
1. Run `npm run build` to build the project
2. Serve the build directory locally
3. Submit forms to test functionality
4. Check browser network tab for successful submissions

### Production Testing
1. Deploy to Netlify
2. Submit test forms
3. Check Netlify Forms dashboard
4. Verify email delivery
5. Test spam protection

### Form Debugging
1. Check Netlify Forms dashboard for submissions
2. Look for form detection in build logs
3. Verify form names match between React and static HTML
4. Ensure all fields are captured correctly

## 📊 Monitoring & Analytics

### Netlify Dashboard
- Form submissions count
- Spam detections
- Field validation errors
- Email delivery status

### Email Delivery
- Monitor `info@cloudpartner.pro` inbox
- Check spam/junk folders initially
- Set up email rules for automatic organization

## 🔒 Security Features

### Implemented Protections
- **Honeypot fields** for bot detection
- **reCAPTCHA** (can be enabled in Netlify dashboard)
- **Rate limiting** (handled by Netlify)
- **Form validation** on client and server side

### Headers Configuration
Security headers are configured in `netlify.toml`:
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- X-Content-Type-Options: nosniff
- Referrer-Policy: strict-origin-when-cross-origin

## 🚨 Troubleshooting

### Forms Not Submitting
1. Check form `name` attribute matches static HTML
2. Verify `data-netlify="true"` is present
3. Ensure hidden `form-name` field exists
4. Check network tab for 200 response

### Emails Not Receiving
1. Check Netlify Forms dashboard for submissions
2. Verify email address in `netlify.toml`
3. Check spam/junk folders
4. Contact Netlify support if persistent

### Template Not Loading
1. Verify template path in `netlify.toml`
2. Check template file exists and is valid HTML
3. Ensure field names match form and template
4. Test with simple template first

### Build Failures
1. Check static HTML files exist in `/public/`
2. Verify form names are consistent
3. Ensure all React form fields have `name` attributes
4. Check build logs for specific errors

## 📝 Maintenance

### Regular Tasks
- Monitor form submissions weekly
- Update email templates for seasonal campaigns
- Test forms after major site updates
- Review spam submissions and adjust filters

### Updating Email Templates
1. Edit templates in `netlify/emails/`
2. Test locally if possible
3. Deploy changes
4. Send test submission to verify

### Adding New Forms
1. Create React form component with Netlify attributes
2. Add static HTML version to `/public/`
3. Create email templates (HTML and text)
4. Update templates in Netlify Dashboard
5. Test thoroughly before production

## 🚀 Deployment & Setup Checklist

### Initial Deployment
1. **Deploy to Netlify** - Forms will be automatically detected during build
2. **Configure Email Templates** - Set up custom email notifications in Netlify Dashboard (see configuration section above)
3. **Test Form Submissions** - Submit test data to verify email delivery
4. **Monitor Email Delivery** - Check `info@cloudpartner.pro` inbox
5. **Configure Spam Protection** - Enable reCAPTCHA if needed via Netlify dashboard

### Post-Deployment Verification
- [ ] Both forms appear in Netlify Dashboard Forms section
- [ ] Email notifications configured for both forms
- [ ] Test submissions received at `info@cloudpartner.pro`
- [ ] Email templates render correctly
- [ ] Spam protection is working
- [ ] Form validation functions properly

## 🔗 Useful Links

- [Netlify Forms Documentation](https://docs.netlify.com/forms/setup/)
- [Netlify Email Templates](https://docs.netlify.com/forms/notifications/)
- [Form Spam Filtering](https://docs.netlify.com/forms/spam-filters/)
- [Netlify Dashboard](https://app.netlify.com/) (Forms section)

## 📞 Support

For technical issues:
1. Check this documentation first
2. Review Netlify Forms dashboard
3. Check build logs for errors
4. Contact development team if needed

For email delivery issues:
1. Check inbox and spam folders
2. Verify email address configuration
3. Contact hosting provider if persistent
4. Consider backup email notification system

---

**Last Updated:** Migration from EmailJS completed
**Next Review:** Monitor for 30 days post-migration