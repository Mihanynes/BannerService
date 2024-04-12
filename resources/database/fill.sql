-- Вставляем новые данные в таблицу features
INSERT INTO features (id, name)
VALUES
    (1, 'Homepage Carousel'),
    (2, 'Search Functionality'),
    (3, 'Customer Support'),
    (4, 'Wishlist Feature'),
    (5, 'Newsletter Subscription'),
    (6, 'Social Media Integration'),
    (7, 'Gift Card Options'),
    (8, 'Multi-language Support'),
    (9, 'Product Comparison Tool'),
    (10, 'Live Chat Support');

-- Вставляем новые данные в таблицу tags
INSERT INTO tags (id, name)
VALUES
    (1, 'Health'),
    (2, 'Entertainment'),
    (3, 'Education'),
    (4, 'Finance');

-- Вставляем новые данные в таблицу user_banners
INSERT INTO user_banners (id, content, is_active, feature_id)
VALUES
    (1, '{
  "title": "Explore Our Featured Products",
  "text": "Check out our handpicked selection of top products.",
  "api": "https://example.com/featured-products"
}', true, 1),
    (2, '{
  "title": "Find Anything You Need",
  "text": "Our powerful search engine helps you discover products quickly.",
  "api": "https://example.com/product-search"
}', true, 2),
    (3, '{
  "title": "24/7 Customer Support",
  "text": "Our support team is here to assist you anytime, anywhere.",
  "api": "https://example.com/customer-support"
}', true, 3),
    (4, '{
  "title": "Create Your Wishlist",
  "text": "Save your favorite items and share them with friends and family.",
  "api": "https://example.com/wishlist"
}', true, 4),
    (5, '{
  "title": "Subscribe to Our Newsletter",
  "text": "Stay updated with the latest news, offers, and promotions.",
  "api": "https://example.com/newsletter-subscription"
}', true, 5),
    (6, '{
  "title": "Connect with Us on Social Media",
  "text": "Follow us on social media for exclusive updates and contests.",
  "api": "https://example.com/social-media"
}', true, 6),
    (7, '{
  "title": "Give the Gift of Choice",
  "text": "Purchase a gift card for your loved ones to shop their favorites.",
  "api": "https://example.com/gift-cards"
}', true, 7),
    (8, '{
  "title": "Choose Your Language",
  "text": "Our website is available in multiple languages for your convenience.",
  "api": "https://example.com/language-options"
}', true, 8),
    (9, '{
  "title": "Compare Products Side by Side",
  "text": "Make informed decisions with our product comparison tool.",
  "api": "https://example.com/product-comparison"
}', true, 9),
    (10, '{
  "title": "Chat with Us Live",
  "text": "Get instant assistance with our live chat support feature.",
  "api": "https://example.com/live-chat"
}', true, 10);

-- Вставляем новые данные в таблицу user_banners_tags
INSERT INTO user_banners_tags (banner_id, tag_id) VALUES
                                                      (1, 1),
                                                      (1, 2),
                                                      (2, 2),
                                                      (2, 3),
                                                      (3, 3),
                                                      (3, 4),
                                                      (4, 4),
                                                      (4, 2),
                                                      (5, 3),
                                                      (5, 1),
                                                      (6, 1),
                                                      (6, 3),
                                                      (7, 2),
                                                      (7, 4),
                                                      (8, 1),
                                                      (8, 4),
                                                      (9, 1),
                                                      (9, 1),
                                                      (10, 3),
                                                      (10, 2);
