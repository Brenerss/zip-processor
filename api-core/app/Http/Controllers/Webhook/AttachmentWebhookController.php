<?php

declare(strict_types = 1);

namespace App\Http\Controllers\Webhook;

use App\Http\Controllers\Controller;
use App\Models\Attachment;
use Illuminate\Http\Request;
use Log;

class AttachmentWebhookController extends Controller
{
    public function update(Request $request, Attachment $attachment)
    {
        Log::info("testing", compact($attachment));
    }
}
