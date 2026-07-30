<?php

declare(strict_types = 1);

namespace App\Http\Controllers\V1;

use App\Http\Controllers\Controller;
use App\Http\Requests\StoreAttachmentRequest;
use App\Services\AttachmentService;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Log;

class AttachmentController extends Controller
{
    public function store(StoreAttachmentRequest $request): Response
    {
        Log::info('AttachmentController@store called with file: ' . $request->file('file')->getClientOriginalName());
        (new AttachmentService($request->file('file')))->store();

        return response()->noContent();
    }
}
