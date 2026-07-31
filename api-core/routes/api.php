<?php

declare(strict_types = 1);

use App\Http\Controllers\V1;
use App\Http\Controllers\Webhook\AttachmentWebhookController;
use Illuminate\Support\Facades\Route;

Route::prefix('/v1')
    ->name('api.v1.')
    ->group(function () {
        Route::post('/auth/test', [V1\AuthController::class, 'loginTestUser'])
            ->name('auth.test');

        Route::middleware('auth:sanctum')
            ->group(function () {
                Route::post('/attachments', [V1\AttachmentController::class, 'store'])
                    ->name('attachments.store');
            });
    });

Route::prefix('/webhooks')
    ->name('webhooks.')
    ->middleware('verify.processor.webhook')
    ->group(function () {
        Route::post('/attachments/{attachment}/status', [AttachmentWebhookController::class, 'update']);
    });
