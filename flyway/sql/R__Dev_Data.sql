-- https://documentation.red-gate.com/fd/repeatable-migrations-273973335.html

-- Only execute in development environment
DO
$$
    BEGIN
        IF '${environment}' = 'dev' THEN
            IF NOT EXISTS (SELECT 1 FROM template_service.template WHERE name = 'Test Template') THEN
                INSERT INTO template_service.template (id,
                                                       name,
                                                       category_id,
                                                       content,
                                                       format,
                                                       version,
                                                       is_active,
                                                       created_by)
                VALUES (uuid_generate_v4(),
                        'Test Template',
                        (SELECT id FROM template_service.template_category WHERE name = 'Email'),
                        '<html>
                <head>
                  <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
                </head>
                <body class="bg-gray-50 font-sans">
                  <div class="max-w-2xl mx-auto my-10 bg-white rounded-lg shadow-lg overflow-hidden">
                    <div class="bg-blue-600 text-white px-6 py-4">
                      <h1 class="text-2xl font-bold">Hello, {{.name}}!</h1>
                    </div>
                    <div class="px-6 py-8">
                      <p class="text-gray-700 mb-6">
                        This is a <span class="font-semibold">test template</span>. Feel free to modify.
                      </p>
                      <div class="bg-blue-50 rounded-lg p-4 border border-blue-200">
                        <h2 class="text-xl font-semibold text-blue-800 mb-3">Details:</h2>
                        <ul class="space-y-2 text-gray-700">
                          <li class="flex"><span class="font-medium w-32">Name:</span> {{.name}}</li>
                          <li class="flex"><span class="font-medium w-32">Date:</span> {{.date}}</li>
                          <li class="flex"><span class="font-medium w-32">Message:</span> {{.message}}</li>
                        </ul>
                      </div>
                      </ul>
                    </div>
                  </div>
                </body>
                </html>',
                        'html',
                        1,
                        true,
                        'system');

                -- Insert template variables
                WITH template_uuid AS (SELECT id
                                       FROM template_service.template
                                       WHERE name = 'Test Template')
                INSERT
                INTO template_service.template_variable (template_id,
                                                         variable_name,
                                                         description,
                                                         default_value,
                                                         is_required)
                VALUES ((SELECT id FROM template_uuid), 'name', 'Recipient name', 'User', true),
                       ((SELECT id FROM template_uuid), 'date', 'Current date', CURRENT_DATE::text, false),
                       ((SELECT id FROM template_uuid), 'message', 'Custom greeting message',
                        'Default message', false);
            END IF;

            -- Add invoice template if it doesn't exist
            IF NOT EXISTS (SELECT 1 FROM template_service.template WHERE name = 'Invoice Template') THEN
                -- Insert invoice template (content abbreviated for brevity)
                INSERT INTO template_service.template (id,
                                                       name,
                                                       category_id,
                                                       content,
                                                       format,
                                                       version,
                                                       is_active,
                                                       created_by)
                VALUES (uuid_generate_v4(),
                        'Invoice Template',
                        (SELECT id FROM template_service.template_category WHERE name = 'Report'),
                        '<html>
                <head>
                  <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
                </head>
                <body class="bg-gray-50 font-sans p-8">
                  <!-- Invoice template content here -->
                  <h1>Invoice for {{.customer_name}}</h1>
                  <p>Invoice #{{.invoice_number}} dated {{.date}}</p>
                </body>
                </html>',
                        'html',
                        1,
                        true,
                        'system');
            END IF;
            IF NOT EXISTS (SELECT 1 FROM template_service.configuration WHERE config_key = 'test_mode') THEN
                INSERT INTO template_service.configuration (config_key, config_value, description)
                VALUES ('test_mode', 'true', 'Enable test mode in development');
            END IF;
        END IF;
    END;
$$;