<?php

declare(strict_types = 1);

namespace App\Http\Controllers\V1;

use App\Http\Controllers\Controller;
use App\Http\Requests\StoreAttachmentRequest;
use App\Services\AttachmentService;
use Illuminate\Http\Response;

class AttachmentController extends Controller
{
    public function store(StoreAttachmentRequest $request): Response
    {
        (new AttachmentService($request->file('file')))->store();

        return response()->noContent();
    }
}
