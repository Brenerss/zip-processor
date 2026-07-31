<?php

declare(strict_types = 1);

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class ProcessorWebhookSignature
{
    /**
     * Handle an incoming request.
     *
     * @param  Closure(Request): (Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        $signature = $request->header('X-Webhook-Signature');
        $timestamp = $request->header('X-Webhook-Timestamp');

        if (!$signature || !$timestamp) {
            return response()->json(['message' => 'Missing webhook headers.'], Response::HTTP_UNAUTHORIZED);
        }

        if (!ctype_digit($timestamp)) {
            return response()->json(['message' => 'Invalid timestamp.'], Response::HTTP_UNAUTHORIZED);
        }

        if (abs(time() - (int) $timestamp) > 300) {
            return response()->json(['message' => 'Webhook expired.'], Response::HTTP_UNAUTHORIZED);
        }

        $payload           = $request->getContent();
        $expectedSignature = hash_hmac('sha256', "$timestamp.$payload", config('services.processor-image.secret'));

        if (!hash_equals($expectedSignature, $signature)) {
            return response()->json(['message' => 'Invalid webhook signature.'], Response::HTTP_UNAUTHORIZED);
        }

        return $next($request);
    }
}
