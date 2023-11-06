package ar.com.montepagano.search.examples.creditcardpayments;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.Scanner;

import ar.com.montepagano.search.middleware.v1.Middleware;
import org.json.JSONObject;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import com.google.protobuf.ByteString;
import ar.com.montepagano.search.middleware.v1.PrivateMiddlewareServiceGrpc;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelRequest;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelResponse;
import ar.com.montepagano.search.appmessage.v1.AppMessageOuterClass.AppSendRequest;
import ar.com.montepagano.search.appmessage.v1.AppMessageOuterClass.AppMessage;
import ar.com.montepagano.search.contracts.v1.Contracts.GlobalContract;
import ar.com.montepagano.search.contracts.v1.Contracts.GlobalContractFormat;

public class Main {
    public static void main(String[] args) {
        System.out.println( "Hello World!" );
        Scanner scanner = new Scanner(System.in);

        // List of book titles
        String[] bookTitles = {
                "The Great Gatsby",
                "To Kill a Mockingbird",
                "1984",
                "Pride and Prejudice",
                "Animal Farm",
                "The Catcher in the Rye",
                "Lord of the Flies",
                "Brave New World",
                "The Hobbit",
                "The Fellowship of the Ring"
        };

        // Prompt user to select book titles
        System.out.println("Please select at least one book to purchase:");
        for (int i = 0; i < bookTitles.length; i++) {
            System.out.println((i + 1) + ". " + bookTitles[i]);
        }

        // Get user input for selected book titles
        System.out.print("Enter the number(s) of the book(s) you want to purchase (comma-separated): ");
        String selectedBooksInput = scanner.nextLine();
        String[] selectedBooksInputArray = selectedBooksInput.split(",");
        int[] selectedBooks = new int[selectedBooksInputArray.length];
        for (int i = 0; i < selectedBooksInputArray.length; i++) {
            selectedBooks[i] = Integer.parseInt(selectedBooksInputArray[i].trim());
        }

        // Prompt user for shipping address
        System.out.print("Enter your shipping address: ");
        String shippingAddress = scanner.nextLine();

        // Print purchase information
        System.out.println("Purchase Information:");
        System.out.println("Selected Books:");
        for (int i = 0; i < selectedBooks.length; i++) {
            System.out.println("- " + bookTitles[selectedBooks[i] - 1]);
        }
        System.out.println("Shipping Address: " + shippingAddress);

        // Load file contract.fsa into a GlobalContract
        ByteString contractBytes = null;
        try {
            contractBytes = ByteString.readFrom(new FileInputStream("./contract.fsa"));
        } catch (IOException e) {
            e.printStackTrace();
        }
        GlobalContract contract = GlobalContract.newBuilder().setContract(contractBytes).setFormat(
                GlobalContractFormat.GLOBAL_CONTRACT_FORMAT_FSA).setInitiatorName("ClientApp").build();

        // Register channel with middleware using GlobalContract
        ManagedChannel channel = ManagedChannelBuilder.forTarget("localhost:1234").usePlaintext().build();
        PrivateMiddlewareServiceGrpc.PrivateMiddlewareServiceBlockingStub stub = PrivateMiddlewareServiceGrpc.newBlockingStub(channel);
        RegisterChannelRequest request = RegisterChannelRequest.newBuilder().setRequirementsContract(contract).build();
        RegisterChannelResponse response = stub.registerChannel(request);
        var channelId = response.getChannelId();

        // Send PurchaseRequest with each item quantities and the shipping address
        var body = ByteString.copyFromUtf8("{}");
        var msg = AppMessage.newBuilder().setType("PurchaseRequest").setBody(body).build();
        var sendreq = AppSendRequest.newBuilder().setChannelId(channelId).setRecipient("Srv").setMessage(msg).build();
        var sendresp = stub.appSend(sendreq);
        if (sendresp.getResult() != Middleware.AppSendResponse.Result.RESULT_OK) {
            System.out.println("Error sending PurchaseRequest. Exiting...");
            System.exit(1);
        }


        channel.shutdown();
    }
}